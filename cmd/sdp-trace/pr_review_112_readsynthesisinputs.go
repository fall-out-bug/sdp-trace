package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

func readPRReviewSynthesisInputs(opts *flagSet) (prReviewSynthesisInputs, error) {
	packet, err := prreview.ReadPacket(opts.stringValue("packet"))
	if err != nil {
		return prReviewSynthesisInputs{}, err
	}
	// Packet, run set, and optional prior ledger are read before synthesis so
	// the output is derived from complete local evidence.
	runs, err := prreview.ReadRunSet(opts.stringValue("runs"))
	if err != nil {
		return prReviewSynthesisInputs{}, err
	}
	existing, err := readOptionalPRReviewLedger(opts.stringValue("existing-ledger"))
	if err != nil {
		return prReviewSynthesisInputs{}, err
	}
	// Existing ledger, when supplied, is an input to synthesis rather than an
	// authority overriding fresh run outputs.
	return prReviewSynthesisInputs{packet: packet, runs: runs, existing: existing}, nil
}
