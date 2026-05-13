package harnessobs

func findUnsafeMapAt(path string, values map[string]any, rawEvent bool) (string, string) {
	// findUnsafeMapAt keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for key, child := range values {

		if field, reason := findUnsafeMapChild(path, key, child, rawEvent); field != "" {
			return field, reason
		}
	}
	return "", ""
}
