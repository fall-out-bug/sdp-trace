package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	// CLI parsing is the only place that translates user input into process
	// exit codes; parsing and joining stay reusable for focused tests.
	opts, err := parseFlags(args, stderr)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}

	rows, err := loadRows(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	// Result rendering owns threshold exit semantics after both inputs are joined.
	return printResults(stdout, stderr, rows, opts.threshold, opts.strictLess)
}

func parseFlags(args []string, output io.Writer) (options, error) {
	flags, opts := newFlagSet(output)
	// Unknown or malformed flags are returned to run so the command can report
	// the parse error and use the documented usage-failure exit code.
	if err := flags.Parse(args); err != nil {
		return options{}, err
	}
	// Both inputs are required because CRAP joins independent coverage and
	// complexity reports by source location.
	if missingInputPath(*opts) {
		return options{}, fmt.Errorf("usage: crapcheck -cover-func cover.txt -gocyclo cyclo.txt [-threshold 5]")
	}
	// No positional arguments are interpreted here; package selection belongs
	// to the tools that produced the input reports.
	return *opts, nil
}
