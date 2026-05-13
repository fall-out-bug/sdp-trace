package harnessobs

func timestampForKey(raw map[string]any, key string) string {
	// timestampForKey keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if observedAt := stringTimestampForKey(raw, key); observedAt != "" {
		return observedAt
	}

	if value, ok := findNumberByKey(raw, key); ok {
		return unixMillisTimestamp(value)
	}
	return ""
}
