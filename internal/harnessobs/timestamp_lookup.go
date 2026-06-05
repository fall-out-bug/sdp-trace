package harnessobs

// Timestamp lookup tries the accepted observation timestamp keys in priority
// order and lets string parsing win over numeric fallback conversion.
func findTimestamp(raw map[string]any) string {
	for _, key := range []string{"time", "timestamp", "created_at", "observed_at"} {
		if observedAt := timestampForKey(raw, key); observedAt != "" {
			return observedAt
		}
	}
	return ""
}

func timestampForKey(raw map[string]any, key string) string {
	if observedAt := stringTimestampForKey(raw, key); observedAt != "" {
		return observedAt
	}
	if value, ok := findNumberByKey(raw, key); ok {
		return unixMillisTimestamp(value)
	}
	return ""
}
