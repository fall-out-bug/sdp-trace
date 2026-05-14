package main

import (
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
