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

func prReviewPacketOptions(opts *flagSet, args []string, outDir string) prreview.PacketOptions {
	// The packet directory is caller-selected so decomposed and combined review
	// flows can publish different artifacts without changing packet identity.
	options := prreview.PacketOptions{OutDir: outDir}
	fillPRReviewPacketIdentity(&options, opts)
	fillPRReviewPacketEvidence(&options, opts, args)
	return options
}

func fillPRReviewPacketIdentity(options *prreview.PacketOptions, opts *flagSet) {
	// Repository and change refs are provenance anchors; they identify the
	// review subject without reading ambient git state.
	options.RepoID = opts.stringValue("repo-id")
	options.ChangeRef = opts.stringValue("change-ref")
	// Commit and diff inputs bind the packet to immutable source facts supplied
	// by the caller.
	options.BaseCommit = opts.stringValue("base")
	options.HeadCommit = opts.stringValue("head")
	options.DiffPath = opts.stringValue("diff")
	// Metadata remains optional context, not authority for the review verdict.
	options.MetadataPath = opts.stringValue("metadata")
}

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
