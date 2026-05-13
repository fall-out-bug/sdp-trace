package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

func buildPRReviewPacket(opts *flagSet, args []string) (prreview.Packet, error) {
	// Repeated context and verification flags remain ordered packet inputs;
	// comma expansion is CLI sugar, not a separate evidence source.
	// BuildPacket owns validation and persistence of the packet directory.
	// CLI parsing only maps flags into the portable packet options shape.
	return prreview.BuildPacket(prReviewPacketOptions(opts, args, opts.stringValue("out")))
}
