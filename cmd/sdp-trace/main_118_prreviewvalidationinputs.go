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
