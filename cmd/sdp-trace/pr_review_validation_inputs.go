package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

type prReviewValidationInputs struct {
	packet  prreview.Packet
	profile prreview.ReviewProfile
	runs    prreview.RunSet
	ledger  prreview.Ledger
}

func readPRReviewValidationInputs(opts *flagSet) (prReviewValidationInputs, error) {
	packet, err := prreview.ReadPacket(opts.stringValue("packet"))
	if err != nil {
		return prReviewValidationInputs{}, err
	}
	// Packet read is separated from profile/run/ledger reads to keep the error
	// boundary precise for callers and tests.
	profile, runs, ledger, err := readPRReviewValidationArtifacts(opts)
	if err != nil {
		return prReviewValidationInputs{}, err
	}
	return prReviewValidationInputs{packet: packet, profile: profile, runs: runs, ledger: ledger}, nil
}

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
