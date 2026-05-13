package harnessobs

func ensureObject(parent map[string]any, key string) map[string]any {
	// ensureObject keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if child, ok := parent[key].(map[string]any); ok {
		return child
	}

	child := map[string]any{}
	parent[key] = child
	return child
}
