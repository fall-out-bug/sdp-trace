package harnessobs

func hasKeyIn(value any, wanted map[string]bool) bool {
	// hasKeyIn keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	switch v := value.(type) {
	case map[string]any:
		return hasKeyInMap(v, wanted)
	case []any:
		return hasKeyInSlice(v, wanted)
	}

	return false
}
