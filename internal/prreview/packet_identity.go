package prreview

import (
	"fmt"
	"time"
)

func newPacket(opts PacketOptions, refs packetRefs, now time.Time, createdBy, ciState string) Packet {
	// newPacket keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	packet := newPacketIdentity(opts)
	attachPacketRefs(&packet, refs)
	attachPacketProvenance(&packet, now, createdBy, ciState)
	packet.UnavailableFields = unavailablePacketFields(opts)
	return packet
}

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

func attachPacketRefs(packet *Packet, refs packetRefs) {
	// attachPacketRefs keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	packet.DiffRef = refs.diff
	packet.MetadataRef = refs.metadata
	packet.ContextRefs = refs.context
	packet.VerificationRefs = refs.verification
}

func attachPacketProvenance(packet *Packet, now time.Time, createdBy, ciState string) {
	// attachPacketProvenance keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	packet.CIState = ciState
	packet.CreatedAt = now.Format(time.RFC3339)
	packet.CreatedBy = createdBy
}
