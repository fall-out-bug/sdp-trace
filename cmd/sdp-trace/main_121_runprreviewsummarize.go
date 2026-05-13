package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

func runPRReviewSummarize(args []string, stdout, stderr io.Writer) int {
	// Summaries are UX copies of validation and ledger state, not new proof.
	opts, code, ok := parsePRReviewSummarizeArgs(args, stderr)
	if !ok {
		return code
	}
	// Summary text is rendered from validation plus ledger evidence only; it is
	// not an independent approval source.
	validation, ledger, err := readPRReviewSummaryInputs(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	summary := prreview.Summarize(validation, ledger)
	if code, ok := writeOptionalPRReviewSummary(opts.stringValue("out"), summary, stderr); !ok {
		return code
	}
	// Stdout mirrors the summary even when a durable summary file is requested.
	// Human-readable summary text does not replace validation or ledger JSON.
	fmt.Fprint(stdout, summary)
	return 0
}
