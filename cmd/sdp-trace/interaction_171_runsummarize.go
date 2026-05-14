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
