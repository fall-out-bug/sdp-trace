package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

func readExistingPRReviewLedger(path string) (*prreview.Ledger, error) {
	ledger, err := prreview.ReadLedger(path)
	if err != nil {
		return nil, err
	}
	// Return a pointer only after the ledger is fully decoded.
	return &ledger, nil
}
