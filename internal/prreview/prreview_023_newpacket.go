package prreview

import (
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
