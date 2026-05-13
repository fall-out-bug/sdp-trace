package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

func readPRReviewPacketAndProfileValues(opts *flagSet) (prreview.Packet, prreview.ReviewProfile, error) {
	packet, err := prreview.ReadPacket(opts.stringValue("packet"))
	if err != nil {
		return prreview.Packet{}, prreview.ReviewProfile{}, err
	}
	// Profile is loaded after packet so packet path failures stay first.
	profile, err := prreview.ReadProfile(opts.stringValue("profile"))
	if err != nil {
		return prreview.Packet{}, prreview.ReviewProfile{}, err
	}
	return packet, profile, nil
}
