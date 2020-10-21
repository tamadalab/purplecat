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

type options struct {
	dest     string
	context  *purplecat.Context
	helpFlag bool
	args     []string
}

func (opts *options) destination() (*os.File, error) {
	if opts.dest == "" {
		return os.Stdout, nil
	}
	dest, err := os.OpenFile(opts.dest, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return dest, nil
}

func (opts *options) isHelpFlag() bool {
	return len(opts.args) == 0 || opts.helpFlag
}

func helpMessage(progName string) string {
	name := filepath.Base(progName)
	return fmt.Sprintf(`%s version %s
%s [OPTIONS] <PROJECTs...|BUILD_FILEs...>
OPTIONS
    -d, --depth <DEPTH>       specifies the depth for parsing (default: 1)
    -f, --format <FORMAT>     specifies the format of the result. Default is 'markdown'.
                              Available values are: CSV, JSON, YAML, XML, and Markdown.
    -l, --level <LOGLEVEL>    specifies the log level. (default: WARN).
                              Available values are: DEBUG, INFO, WARN, and FATAL
    -o, --output <FILE>       specifies the destination file (default: STDOUT).
    -N, --offline             offline mode (no network access).

    -h, --help                prints this message.
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

func constructFlags(args []string, opts *options, logLevel *string) *flag.FlagSet {
	flags := flag.NewFlagSet("purplecat", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(helpMessage(args[0])) }
	flags.BoolVarP(&opts.context.DenyNetworkAccess, "offline", "N", false, "offline mode (no network access)")
	flags.BoolVarP(&opts.helpFlag, "help", "h", false, "print this message")
	flags.StringVarP(logLevel, "level", "l", "WARN", "specifies the log level")
	flags.IntVarP(&opts.context.Depth, "depth", "d", 1, "specifies the depth for parsing")
	flags.StringVarP(&opts.dest, "output", "o", "", "specifies the destination file (default: STDOUT)")
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

func parseArgs(args []string) (*options, error) {
	opts := &options{context: &purplecat.Context{}}
	var logLevel string
	flags := constructFlags(args, opts, &logLevel)
	if err := flags.Parse(args); err != nil {
		return opts, err
	}
	updateLogLevel(logLevel)
	opts.args = flags.Args()[1:]
	return opts, nil
}

func performEach(projectPath string, context *purplecat.Context) (*purplecat.Project, error) {
	parser, err := context.GenerateParser(projectPath)
	if err != nil {
		return nil, err
	}
	return parser.Parse(purplecat.NewPath(projectPath))
}

func perform(opts *options) int {
	dest, err := opts.destination()
	if err != nil {
		return printError(err, 9)
	}
	writer, err2 := opts.context.NewWriter(dest)
	if err2 != nil {
		return printError(err2, 8)
	}

	for _, project := range opts.args {
		tree, err := performEach(project, opts.context)
		if err != nil {
			return printError(err, 2)
		}
		writer.Write(tree)
	}
	return 0
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
