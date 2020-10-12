package main

import (
	"fmt"
	"io"
	"os"

	flag "github.com/spf13/pflag"
)

const VERSION = "1.0.0"

type options struct {
	dest     string
	offline  bool
	helpFlag bool
	destFile *os.File
	args     []string
}

func (opts *options) finish() {
	if opts.destFile != nil {
		opts.destFile.Close()
	}
}

func (opts *options) destination() (io.Writer, error) {
	if opts.dest == "" {
		return os.Stdout, nil
	}
	dest, err := os.OpenFile(opts.dest, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	opts.destFile = dest
	return dest, nil
}

func (opts *options) isHelpFlag() bool {
	return len(opts.args) == 0 || opts.helpFlag
}

func helpMessage(progName string) string {
	return fmt.Sprintf(`%s version %s
%s [OPTIONS] <PROJECTs...>
OPTIONS
	-d, --dest <FILE>    specifies the destination file (default: STDOUT).
	-N, --offline        offline mode (no network access).

	-h, --help           prints this message.
PROJECT
    target project for extracting related libraries and their licenses.`, progName, VERSION, progName)
}

func printError(err error, status int) int {
	if err != nil {
		fmt.Println(err.Error())
		return status
	}
	return 0
}

func constructFlags(opts *options) *flag.FlagSet {
	flags := flag.NewFlagSet("purplecat", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(helpMessage(args[0])) }
	flags.BoolVarP(&opts.offline, "offline", "N", false, "offline mode (no network access)")
	flags.BoolVarP(&opts.helpFlag, "help", "h", false, "print this message")
	flags.StringVarP(&opts.dest, "dest", "d", "", "specifies the destination file (default: STDOUT)")
	return flags
}

func parseArgs(args []string) (*options, int, error) {
	opts := &options{}
	flags := constructFlags(opts)
	if err := flags.Parse(args); err != nil {
		return opts, 1, err
	}
	if opts.isHelpFlag() {
		return opts, 0, fmt.Errorf(helpMessage(args[0]))
	}
	opts.args = flags.Args()[1:]
	return opts, 0, nil
}

func perform(opts *options) int {
	dest, err := opts.destination()
	if err != nil {
		return 9
	}
	defer opts.finish()
	for _, project := range opts.args {
		tree, err := purplecat.ParseProject(project)
		if err != nil {
			return printError(err, 2)
		}
		tree.Println(dest)
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
