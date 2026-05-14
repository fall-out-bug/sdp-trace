package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

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
