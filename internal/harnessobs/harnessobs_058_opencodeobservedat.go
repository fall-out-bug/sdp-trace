package harnessobs

import (
	"time"
)

func openCodeObservedAt(raw map[string]any, now time.Time) string {
	// openCodeObservedAt keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	observedAt := findTimestamp(raw)
	if observedAt == "" {

		return now.Format(time.RFC3339)
	}
	return observedAt
}
