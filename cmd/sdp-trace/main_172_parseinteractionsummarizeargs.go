package main

import (
	"fmt"
	"io"
)

func parseInteractionSummarizeArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "interaction summarize"}
	// The summary command accepts a trace artifact and optional output only; it
	// does not accept ad hoc report text.
	opts.setString("trace", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if !requireOnlyFlags(opts, stderr, "interaction summarize accepts only flags", interactionSummarizeRequiredFlags) {
		return nil, exitUsage, false
	}
	return opts, 0, true
}
