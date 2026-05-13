package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/harnessobs"
)

func runHarnessSummarize(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "harness summarize"}
	// Summaries are read-only views over a validation artifact.
	opts.setString("validation", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	// Summarize accepts only a persisted validation artifact, not raw events.
	if !requireOnlyFlags(opts, stderr, "harness summarize accepts only flags", harnessSummarizeRequiredFlags) {
		return exitUsage
	}
	validation, err := harnessobs.LoadValidation(opts.stringValue("validation"))
	if err != nil {
		// Unreadable validation artifacts keep summary in cannot_verify because
		// there is no trusted row set to summarize.
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	// Summaries are derived from persisted validation evidence; missing or
	// malformed validation stays cannot_verify instead of becoming prose truth.
	fmt.Fprint(stdout, harnessobs.Summarize(validation))
	return 0
}
