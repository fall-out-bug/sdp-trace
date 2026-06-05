package main

import (
	"fmt"
	"io"
	"os"

	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

func readPRReviewSummaryInputs(opts *flagSet) (prreview.Validation, prreview.Ledger, error) {
	validation, err := prreview.ReadValidation(opts.stringValue("validation"))
	if err != nil {
		return prreview.Validation{}, prreview.Ledger{}, err
	}
	// Ledger loading happens after validation so summary failures identify the
	// missing artifact class without masking the validation path error.
	ledger, err := prreview.ReadLedger(opts.stringValue("ledger"))
	return validation, ledger, err
}

func writeOptionalPRReviewSummary(path, summary string, stderr io.Writer) (int, bool) {
	if path == "" {
		// Summary output is optional; no path means there is no write attempt and
		// no additional failure state.
		return 0, true
	}
	return writePRReviewSummaryFile(path, summary, stderr)
}

func writePRReviewSummaryFile(path, summary string, stderr io.Writer) (int, bool) {
	if err := refuseExistingFile(path); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage, false
	}
	// Summaries are write-once CLI artifacts; refusing overwrite protects the
	// cited review text from accidental replacement.
	if err := os.WriteFile(path, []byte(summary), 0o644); err != nil {
		fmt.Fprintln(stderr, err)
		return 1, false
	}
	return 0, true
}
