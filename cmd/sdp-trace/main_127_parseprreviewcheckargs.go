package main

import (
	"fmt"
	"io"
)

func parsePRReviewCheckArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "pr-review check"}
	registerPRReviewCheckFlags(opts)
	if err := opts.parse(args); err != nil {
		// The combined command still has a flag-only contract.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if rejectRest(opts, stderr, "pr-review check accepts only flags") {
		// Positional payload would bypass the packet's declared provenance.
		return nil, exitUsage, false
	}
	// Required review anchors are checked after parsing so diagnostics reflect
	// the declared command shape.
	if err := requirePRReviewCheckInputs(opts); err != nil {
		// Missing anchors are caught before any reviewer process can run.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	// Successful parsing only validates command shape; execution still has to
	// build packet, profile, and run evidence.
	return opts, 0, true
}
