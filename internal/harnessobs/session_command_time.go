package harnessobs

import "time"

func sessionCommandObservedAt(session SessionRun) string {
	observedAt := session.StartTime
	if observedAt == "" {
		observedAt = session.CreatedAt
	}
	if _, err := time.Parse(time.RFC3339, observedAt); err != nil {
		return time.Now().UTC().Format(time.RFC3339)
	}
	return observedAt
}
