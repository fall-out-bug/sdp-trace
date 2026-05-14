package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

func readOptionalPRReviewLedger(path string) (*prreview.Ledger, error) {
	if path == "" {
		// A missing optional ledger starts synthesis from an empty review record.
		return nil, nil
	}
	return readExistingPRReviewLedger(path)
}
