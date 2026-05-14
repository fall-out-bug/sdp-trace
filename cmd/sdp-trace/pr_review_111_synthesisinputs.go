package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

type prReviewSynthesisInputs struct {
	packet   prreview.Packet
	runs     prreview.RunSet
	existing *prreview.Ledger
}
