package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

func fillPRReviewPacketEvidence(options *prreview.PacketOptions, opts *flagSet, args []string) {
	// Repeated path flags are reconstructed from raw args so order survives the
	// simple parser's single-value storage.
	options.ContextPaths = repeatedFlagValues(args, "context", opts.stringValue("context"))
	options.VerificationPaths = repeatedFlagValues(args, "verification", opts.stringValue("verification"))
	// CI state and producer are declared packet metadata, not inferred local
	// evidence.
	options.CIState = opts.stringValue("ci-state")
	options.CreatedBy = opts.stringValue("created-by")
}
