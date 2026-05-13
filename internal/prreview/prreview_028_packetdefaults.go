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
	createdBy := strings.TrimSpace(opts.CreatedBy)
	if createdBy == "" {
		createdBy = "sdp-trace"
	}
	ciState := opts.CIState
	if ciState == "" {

		ciState = StateNotAssessed
	}
	return now, createdBy, ciState
}
