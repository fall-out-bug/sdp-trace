package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

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
