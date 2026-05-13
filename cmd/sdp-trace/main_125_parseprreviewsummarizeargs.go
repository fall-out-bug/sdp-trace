package main

import (
	"fmt"
	"io"
)

func parsePRReviewSummarizeArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "pr-review summarize"}
	// Summaries accept only evidence paths and an optional output path; any extra
	// positional text would be unaudited report content.
	opts.setString("validation", "")
	opts.setString("ledger", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if rejectRest(opts, stderr, "pr-review summarize accepts only flags") {
		return nil, exitUsage, false
	}
	// Required inputs are loaded by readPRReviewSummaryInputs so bad paths remain
	// cannot_verify rather than usage-only failures.
	return opts, 0, true
}
