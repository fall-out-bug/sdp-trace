package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

func readPRReviewPacketAndProfile(opts *flagSet, stderr io.Writer) (prreview.Packet, prreview.ReviewProfile, bool) {
	packet, profile, err := readPRReviewPacketAndProfileValues(opts)
	if err != nil {
		// Packet/profile load failures mean review evidence cannot be replayed.
		fmt.Fprintln(stderr, err)
		return prreview.Packet{}, prreview.ReviewProfile{}, false
	}
	return packet, profile, true
}
