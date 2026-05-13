package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

func readPRReviewValidationArtifacts(opts *flagSet) (prreview.ReviewProfile, prreview.RunSet, prreview.Ledger, error) {
	profile, err := prreview.ReadProfile(opts.stringValue("profile"))
	if err != nil {
		return prreview.ReviewProfile{}, prreview.RunSet{}, prreview.Ledger{}, err
	}
	// All validation artifacts are local files; missing or malformed rows keep
	// the final review state from being promoted.
	runs, err := prreview.ReadRunSet(opts.stringValue("runs"))
	if err != nil {
		return prreview.ReviewProfile{}, prreview.RunSet{}, prreview.Ledger{}, err
	}
	// Ledger is read last because it is derived from packet and runs.
	ledger, err := prreview.ReadLedger(opts.stringValue("ledger"))
	return profile, runs, ledger, err
}
