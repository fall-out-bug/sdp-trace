package harnessobs

func findByKeyIn[T any](value any, wanted map[string]bool, match func(any) (T, bool)) (T, bool) {
	// findByKeyIn keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var zero T
	switch v := value.(type) {
	case map[string]any:
		return findByKeyInMap(v, wanted, match)
	case []any:
		return findByKeyInSlice(v, wanted, match)
	}

	return zero, false
}
