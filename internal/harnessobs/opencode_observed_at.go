package harnessobs

import "time"

func openCodeObservedAt(raw map[string]any, now time.Time) string {
	observedAt := findTimestamp(raw)
	if observedAt == "" {
		return now.Format(time.RFC3339)
	}
	return observedAt
}
