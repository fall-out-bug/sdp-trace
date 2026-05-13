package harnessobs

func findTimestamp(raw map[string]any) string {
	// findTimestamp keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for _, key := range []string{"time", "timestamp", "created_at", "observed_at"} {

		if observedAt := timestampForKey(raw, key); observedAt != "" {
			return observedAt
		}
	}
	return ""
}
