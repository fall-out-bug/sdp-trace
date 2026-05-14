package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

func prReviewPacketOptions(opts *flagSet, args []string, outDir string) prreview.PacketOptions {
	// The packet directory is caller-selected so decomposed and combined review
	// flows can publish different artifacts without changing packet identity.
	options := prreview.PacketOptions{OutDir: outDir}
	fillPRReviewPacketIdentity(&options, opts)
	fillPRReviewPacketEvidence(&options, opts, args)
	return options
}
