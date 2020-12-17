package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/tamadalab/purplecat"
	"github.com/tamadalab/purplecat/logger"
)

type serverOpts struct {
	runServer bool
	port      int
}

type commonOpts struct {
	cachePath string
	cacheType string
	logLevel  string
	helpFlag  bool
}

type cliOptions struct {
	dest string
	args []string
}

type options struct {
	context *purplecat.Context
	server  *serverOpts
	common  *commonOpts
	cli     *cliOptions
}

func (opts *options) destination() (*os.File, error) {
	if opts.cli.dest == "" {
		return os.Stdout, nil
	}
	dest, err := os.OpenFile(opts.cli.dest, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return dest, nil
}

func (opts *options) isHelpFlag() bool {
	return opts.common.helpFlag || (len(opts.cli.args) == 0 && !opts.server.runServer)
}

func helpMessage(progName string) string {
	name := filepath.Base(progName)
	return fmt.Sprintf(`%s version %s
%s [COMMON_OPTIONS] [CLI_MODE_OPTIONS] [SERVER_MODE_OPTIONS] <PROJECTs...|BUILD_FILEs...>
COMMON_OPTIONS
    -c, --cache-type <TYPE>        specifies the cache type. (default: default).
                                   Available values are: default, ref-only, newdb and memory.
        --cachedb-path <DBPATH>    specifies the cache database path
                                   (default: ~/.config/purplecat/cachedb.json).
    -l, --log-level <LOGLEVEL>     specifies the log level. (default: WARN).
                                   Available values are: DEBUG, INFO, WARN, and FATAL
    -h, --help                     prints this message.

CLI_MODE_OPTIONS
    -d, --depth <DEPTH>            specifies the depth for parsing (default: 1)
    -f, --format <FORMAT>          specifies the result format. Default is 'markdown'.
                                   Available values are: CSV, JSON, YAML, XML, and Markdown.
    -o, --output <FILE>            specifies the destination file (default: STDOUT).
    -N, --offline                  offline mode (no network access).

SERVER_MODE_OPTIONS
    -p, --port <PORT>              specifies the port number of REST API server. Default is 8080.
                                   If '--server' option did not specified, purplecat ignores this option.
    -s, --server                   starts REST API server. With this option, purplecat ignores
                                   CLI_MODE_OPTIONS and arguments.

PROJECT
    target project for extracting dependent libraries and their licenses.
BUILD_FILE
    build file of the project for extracting dependent libraries and their licenses

purplecat support the projects using the following build tools.
    * Maven 3 (pom.xml)`, name, purplecat.Version, name)
}

func printError(err error, status int) int {
	if err != nil {
		logger.Fatal(err.Error())
		return status
	}
	return 0
}

func constructFlags(args []string, opts *options) *flag.FlagSet {
	flags := flag.NewFlagSet("purplecat", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(helpMessage(args[0])) }
	flags.BoolVarP(&opts.context.DenyNetworkAccess, "offline", "N", false, "offline mode (no network access)")
	flags.BoolVarP(&opts.common.helpFlag, "help", "h", false, "print this message")
	flags.StringVarP(&opts.common.cacheType, "cache-type", "c", "default", "specifies the cache type")
	flags.StringVarP(&opts.common.cachePath, "cachedb-path", "", purplecat.DefaultCacheDBPath(), "specifies the cache database path.")
	flags.StringVarP(&opts.common.logLevel, "log-level", "l", "WARN", "specifies the log level")
	flags.IntVarP(&opts.context.Depth, "depth", "d", 1, "specifies the depth for parsing")
	flags.IntVarP(&opts.server.port, "port", "p", 8080, "specifies the port number of REST API server")
	flags.BoolVarP(&opts.server.runServer, "server", "s", false, "starts REST API server")
	flags.StringVarP(&opts.cli.dest, "output", "o", "", "specifies the destination file (default: STDOUT)")
	flags.StringVarP(&opts.context.Format, "format", "f", "markdown", "specifies the result format (default: markdown).")
	return flags
}

func updateLogLevel(level string) {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		logger.SetLevel(logger.DEBUG)
	case "info":
		logger.SetLevel(logger.INFO)
	case "warn":
		logger.SetLevel(logger.WARN)
	case "fatal":
		logger.SetLevel(logger.FATAL)
	}
}

func validateCacheType(opts *options) error {
	return generalValidator([]string{"default", "ref-only", "newdb", "memory"}, opts.common.cacheType, "%s: unknown cache type")
}

func validateCachePath(opts *options) error {
	stat, err := os.Stat(opts.common.cachePath)
	if err != nil || stat.Mode().IsRegular() {
		return nil
	}
	return fmt.Errorf("%s: not regular file", opts.common.cachePath)
}

func validateFormat(opts *options) error {
	return generalValidator([]string{"csv", "json", "markdown", "yaml", "xml"}, opts.context.Format, "%s: unknown format")
}
func validateLogLevel(opts *options) error {
	return generalValidator([]string{"debug", "info", "warn", "fatal"}, opts.common.logLevel, "%s: unknown log level")
}

func generalValidator(available []string, value, message string) error {
	lower := strings.ToLower(value)
	for _, value := range available {
		if value == lower {
			return nil
		}
	}
	return fmt.Errorf(message, value)
}

func validate(opts *options) error {
	validators := [](func(opts *options) error){
		validateCacheType,
		validateCachePath,
		validateFormat,
		validateLogLevel,
	}
	for _, validator := range validators {
		if err := validator(opts); err != nil {
			return err
		}
	}
	return nil
}

func initializeCache(opts *options) (*options, error) {
	cType := purplecat.ParseCacheType(opts.common.cacheType)
	cc, err := purplecat.NewCacheDBWithPath(cType, opts.common.cachePath)
	if err != nil {
		return opts, err
	}
	opts.context.Cache = cc
	return opts, nil
}

func parseArgs(args []string) (*options, error) {
	opts := &options{context: &purplecat.Context{}, server: &serverOpts{}, cli: &cliOptions{}, common: &commonOpts{}}
	flags := constructFlags(args, opts)
	if err := flags.Parse(args); err != nil {
		return opts, err
	}
	if err := validate(opts); err != nil {
		return opts, err
	}
	updateLogLevel(opts.common.logLevel)
	opts.cli.args = flags.Args()[1:]
	return initializeCache(opts)
}

func performEach(projectPath string, context *purplecat.Context) (*purplecat.Project, error) {
	parser, err := context.GenerateParser(projectPath)
	if err != nil {
		return nil, err
	}
	return parser.Parse(purplecat.NewPath(projectPath))
}

func createWriter(opts *options) (purplecat.Writer, error) {
	dest, err := opts.destination()
	if err != nil {
		return nil, err
	}
	writer, err2 := opts.context.NewWriter(dest)
	if err2 != nil {
		return nil, err2
	}
	return writer, nil
}

func postProcess(context *purplecat.Context) int {
	if err := context.Cache.Store(); err != nil {
		return printError(err, 8)
	}
	return 0
}

func performCli(opts *options) int {
	writer, err := createWriter(opts)
	if err != nil {
		return printError(err, 9)
	}
	for _, project := range opts.cli.args {
		tree, err := performEach(project, opts.context)
		if err != nil {
			return printError(err, 2)
		}
		writer.Write(tree)
	}
	return postProcess(opts.context)
}

func perform(opts *options) int {
	if opts.server.runServer {
		return opts.server.StartServer(opts.common, opts.context.Cache)
	}
	return performCli(opts)
}

func goMain(args []string) int {
	opts, err := parseArgs(args)
	if err != nil {
		return printError(err, 1)
	}
	if opts.isHelpFlag() {
		fmt.Println(helpMessage(args[0]))
		return 0
	}
	return perform(opts)
}

func main() {
	status := goMain(os.Args)
	os.Exit(status)
}

func init() {
	logger.SetLevel(logger.WARN)
}
