package prreview

import (
	"strings"
	"time"
)

func packetDefaults(opts PacketOptions) (time.Time, string, string) {
	// packetDefaults keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	now := opts.Now
	if now.IsZero() {

		now = time.Now().UTC()
	}
	return now, defaultPacketCreator(opts.CreatedBy), defaultPacketCIState(opts.CIState)
}

func defaultPacketCreator(createdBy string) string {
	createdBy = strings.TrimSpace(createdBy)
	if createdBy == "" {
		return "sdp-trace"
	}
	return createdBy
}

func defaultPacketCIState(ciState string) string {
	if ciState == "" {
		return StateNotAssessed
	}
	return ciState
}
