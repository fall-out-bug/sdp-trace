package prreview

import (
	"time"
)

func attachPacketProvenance(packet *Packet, now time.Time, createdBy, ciState string) {
	// attachPacketProvenance keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	packet.CIState = ciState
	packet.CreatedAt = now.Format(time.RFC3339)
	packet.CreatedBy = createdBy
}
