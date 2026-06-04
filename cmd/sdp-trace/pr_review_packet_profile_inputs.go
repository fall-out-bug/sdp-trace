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

func readPRReviewPacketAndProfileValues(opts *flagSet) (prreview.Packet, prreview.ReviewProfile, error) {
	packet, err := prreview.ReadPacket(opts.stringValue("packet"))
	if err != nil {
		return prreview.Packet{}, prreview.ReviewProfile{}, err
	}
	// Profile is loaded after packet so packet path failures stay first.
	profile, err := prreview.ReadProfile(opts.stringValue("profile"))
	if err != nil {
		// Profile errors return an empty packet so callers do not mix partial inputs.
		return prreview.Packet{}, prreview.ReviewProfile{}, err
	}
	return packet, profile, nil
}
