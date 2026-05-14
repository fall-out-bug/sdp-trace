package main

import (
	"fmt"
	"io"
)

func parseCheckpointVerifyArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := newCheckpointVerifyFlagSet()
	if err := opts.parse(args); err != nil {
		// Parse errors happen before any signed checkpoint or policy is loaded.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if len(opts.rest()) != 0 {
		// Positional arguments would make the verification source ambiguous.
		fmt.Fprintln(stderr, "checkpoint verify accepts only flags")
		return nil, exitUsage, false
	}
	if err := requireCheckpointVerifyInputs(opts); err != nil {
		// Required verification inputs identify the run replay source and signed
		// checkpoint artifact before optional policy can affect trust scope.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	return opts, 0, true
}
