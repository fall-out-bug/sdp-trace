package posture

import (
	"fmt"
	"strings"
	"time"
)

func parseFreshnessBoundary(value string) (time.Time, bool, error) {
	// parseFreshnessBoundary keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	if strings.TrimSpace(value) == "" {
		return time.Time{}, false, nil
	}
	if strings.HasPrefix(value, "P") {

		return time.Time{}, false, fmt.Errorf("duration freshness boundaries are not supported in v1")
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, false, err
	}
	return parsed, true, nil
}

// isStale applies the temporal evidence boundary. Parse failures default to stale (safe).
// InputObservedAt before cutoff crosses the freshness threshold.
func isStale(value string, cutoff time.Time) bool {
	// isStale keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return true
	}
	return parsed.Before(cutoff)
}
