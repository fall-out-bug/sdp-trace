package recorder

import (
	"time"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func (w *runWriter) newEvent(eventType trace.EventType, payloadMap map[string]any) trace.Event {
	// Event defaults are applied after recorder-owned provenance fields are set
	// so schema defaults cannot erase the recorded chain linkage.
	return trace.Event{
		SchemaVersion: trace.SchemaVersion,
		RunID:         w.manifest.RunID,

		EventID:   randomHex(24),
		Sequence:  w.sequence,
		EventType: eventType,
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),

		PrevEventHash: w.lastHash,
		HashAlgorithm: trace.HashAlgSHA256,

		Canonicalization: canonicalization(),
		EventPayload:     payloadMap,
		ObservedBy:       "local_recorder",
	}.EnsureDefaults()
}

func canonicalization() trace.Canonicalization {
	// The recorder emits the canonicalization contract alongside every event so
	// offline replay can choose the same hashing algorithm.
	return trace.Canonicalization{
		Algorithm: trace.CanonicalSchemaAlgo,
		Version:   trace.CanonicalAlgoV,
	}
}
