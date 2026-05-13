package harnessobs

import (
	"time"
)

func sessionCommandObservedAt(session SessionRun) string {
	// sessionCommandObservedAt keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	observedAt := session.StartTime
	if observedAt == "" {
		observedAt = session.CreatedAt
	}
	if _, err := time.Parse(time.RFC3339, observedAt); err != nil {

		return time.Now().UTC().Format(time.RFC3339)
	}
	return observedAt
}
