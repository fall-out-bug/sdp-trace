package prreview

import (
	"fmt"
)

func newPacketIdentity(opts PacketOptions) Packet {
	// newPacketIdentity keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	return Packet{
		SchemaVersion:  SchemaVersionPacket,
		PacketID:       fmt.Sprintf("%s-%s-%s", opts.RepoID, opts.ChangeRef, opts.HeadCommit[:12]),
		RepoID:         opts.RepoID,
		ChangeRef:      opts.ChangeRef,
		BaseCommit:     opts.BaseCommit,
		HeadCommit:     opts.HeadCommit,
		RedactionState: RedactionNone,
	}
}
