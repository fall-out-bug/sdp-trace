package main

import (
	"fmt"
	"io"
)

var runHarnessObserveCommand = func(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parseHarnessObserveArgs(args, stderr)
	if !ok {
		return code
	}
	run, err := observeHarnessRun(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	return writeHarnessRun(stdout, stderr, run)
}
