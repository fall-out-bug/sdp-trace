package main

import (
	"fmt"
	"io"
)

func requireWrappedCommandArgs(opts *flagSet, command []string, stderr io.Writer) bool {
	if len(command) == 0 {
		// Recorder runs require an observed command to produce meaningful trace
		// evidence.
		fmt.Fprintln(stderr, "run requires a command")
		return false
	}
	if opts.stringValue("task") == "" {
		// Task id binds the run to the SpecKit task vocabulary.
		fmt.Fprintln(stderr, "run requires --task")
		return false
	}
	if missingRequiredContract(opts) {
		// Contract choice must be explicit unless the caller opts into the
		// built-in default.
		fmt.Fprintln(stderr, "run requires --contract unless --use-default-contract is set")
		return false
	}
	return true
}
