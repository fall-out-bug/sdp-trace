package trace

import (
	"fmt"
	"time"
)

// This file owns construction of newly appended run events.

func newAppendedRunEvent(artifact RunArtifact, eventType EventType, payload map[string]any, observedBy string) (Event, error) {
	// One timestamp feeds both event identity and event time to avoid avoidable
	// intra-event clock skew.
	now := time.Now().UTC()
	sequence := len(artifact.Events)
	event := Event{
		SchemaVersion: SchemaVersion,
		RunID:         artifact.Manifest.RunID,
		EventID:       SHA256Hex(fmt.Sprintf("%s:%s:%d:%s", artifact.Manifest.RunID, eventType, sequence, now.Format(time.RFC3339Nano))),
		Sequence:      sequence,
		EventType:     eventType,
		Timestamp:     now.Format(time.RFC3339Nano),
		PrevEventHash: appendedPrevHash(artifact.Events),
		HashAlgorithm: HashAlgSHA256,
		Canonicalization: Canonicalization{
			Algorithm: CanonicalSchemaAlgo,
			Version:   CanonicalAlgoVersion,
		},
		EventPayload: payload,
		ObservedBy:   observedBy,
	}
	return event.WithComputedEventHash()
}

func appendedPrevHash(events []Event) string {
	// The first retained event links to the explicit null sentinel.
	if len(events) == 0 {
		return NullEventHash
	}
	return events[len(events)-1].EventHash
}
