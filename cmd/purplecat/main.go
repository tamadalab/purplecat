package main

import (
	"fmt"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"
	"github.com/tamadalab/purplecat"
)

type options struct {
	dest     string
	offline  bool
	format   string
	depth    int
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
%s [OPTIONS] <PROJECTs...>
OPTIONS
    -d, --depth <DEPTH>      specifies the depth for parsing (default: 1)
    -f, --format <FORMAT>    specifies the format of the result. Default is 'markdown'.
                             Available values are: CSV, JSON, YAML, XML, and Markdown.
    -o, --output <FILE>      specifies the destination file (default: STDOUT).
    -N, --offline            offline mode (no network access).

    -h, --help               prints this message.
PROJECT
    target project for extracting related libraries and their licenses.`, name, purplecat.VERSION, name)
}

func printError(err error, status int) int {
	if err != nil {
		fmt.Println(err.Error())
		return status
	}
	return 0
}

func constructFlags(args []string, opts *options) *flag.FlagSet {
	flags := flag.NewFlagSet("purplecat", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(helpMessage(args[0])) }
	flags.BoolVarP(&opts.offline, "offline", "N", false, "offline mode (no network access)")
	flags.BoolVarP(&opts.helpFlag, "help", "h", false, "print this message")
	flags.IntVarP(&opts.depth, "depth", "d", 1, "specifies the depth for parsing")
	flags.StringVarP(&opts.dest, "output", "o", "", "specifies the destination file (default: STDOUT)")
	flags.StringVarP(&opts.format, "format", "f", "markdown", "specifies the result format (default: markdown).")
	return flags
}

func parseArgs(args []string) (*options, int, error) {
	opts := &options{}
	flags := constructFlags(args, opts)
	if err := flags.Parse(args); err != nil {
		return opts, 1, err
	}
	opts.args = flags.Args()[1:]
	if opts.isHelpFlag() {
		return opts, 0, fmt.Errorf(helpMessage(args[0]))
	}
	return opts, 0, nil
}

func performEach(projectPath string, context *purplecat.Context) (*purplecat.DependencyTree, error) {
	parser, err := context.GenerateParser(projectPath)
	if err != nil {
		return nil, err
	}
	return parser.Parse(projectPath)
}

func perform(opts *options) int {
	context := purplecat.NewContext(!opts.offline, opts.format, opts.depth)
	dest, err := opts.destination()
	if err != nil {
		return printError(err, 9)
	}
	writer, err2 := context.NewWriter(dest)
	if err2 != nil {
		return printError(err2, 8)
	}

	for _, project := range opts.args {
		tree, err := performEach(project, context)
		if err != nil {
			return printError(err, 2)
		}
		writer.Write(tree)
	}
	return 0
}

func goMain(args []string) int {
	opts, status, err := parseArgs(args)
	if err != nil {
		return printError(err, status)
	}
	return perform(opts)
}

func main() {
	status := goMain(os.Args)
	os.Exit(status)
}
