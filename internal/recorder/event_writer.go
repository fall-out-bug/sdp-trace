package recorder

import "github.com/fall_out_bug/sdp-trace/internal/trace"

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
