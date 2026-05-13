package recorder

import (
	"path/filepath"
	"time"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

// Writer event operations own the append-only chain. Each persisted event is
// immediately reflected into the manifest so disk state and in-memory state
// advance together.
//
// The chain update sequence is deliberately local to this file: normalize the
// payload, compute the event hash with the trace package, write the event file,
// advance the in-memory head, then rewrite the manifest. Callers should not be
// able to skip one of those steps.
//
// Event writing remains synchronous so a returned error always means the chain
// or manifest was not advanced for that event.
//
// That invariant is what lets callers treat an append error as an open run
// rather than a partially successful gate verdict.

func (w *runWriter) appendEvent(eventType trace.EventType, payload any) error {
	// Events are converted, hashed, persisted, then reflected into the manifest
	// in that order so the manifest never points at an unwritten event.
	payloadMap, err := toEventPayload(payload)
	if err != nil {
		return err
	}
	event := w.newEvent(eventType, payloadMap)
	computed, err := event.WithComputedEventHash()
	if err != nil {
		return err
	}
	if err := w.persistEvent(computed, eventType); err != nil {
		return err
	}
	w.advanceEventHead(computed)

	return w.writeManifest()
}

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

func (w *runWriter) persistEvent(event trace.Event, eventType trace.EventType) error {
	// Sequence-prefixed filenames keep the append order visible in filesystem
	// listings without making the filename part of the event hash.
	filename := filepath.Join(w.runDir, "events", eventFilename(event.Sequence, eventType))
	return writeIndentedJSON(filename, event)
}

func (w *runWriter) advanceEventHead(event trace.Event) {
	// The in-memory head mirrors the persisted event chain and feeds both the
	// next event's previous hash and the manifest closure fields.
	w.sequence++
	w.lastHash = event.EventHash
	w.events = append(w.events, event)
}
