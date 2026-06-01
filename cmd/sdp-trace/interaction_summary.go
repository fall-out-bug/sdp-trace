package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/interaction"
)

func runInteractionSummarize(args []string, stdout, stderr io.Writer) int {
	// Summarize reads an existing trace artifact and emits a derived report; it
	// never imports or mutates trace events.
	opts, code, ok := parseInteractionSummarizeArgs(args, stderr)
	if !ok {
		return code
	}
	// Summaries are derived views of trace evidence; unreadable traces keep the
	// command from producing an overconfident report.
	trace, err := interaction.ReadTrace(opts.stringValue("trace"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	// Summary generation is pure over the decoded trace.
	summary := interaction.SummarizeTrace(trace)
	if err := writeOptionalJSONFile(opts.stringValue("out"), summary); err != nil {
		// Optional output write failures are publication failures after a valid
		// trace summary was derived.
		fmt.Fprintln(stderr, err)
		return 1
	}
	writeJSONPayloadUnchecked(stdout, summary)
	return 0
}

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
