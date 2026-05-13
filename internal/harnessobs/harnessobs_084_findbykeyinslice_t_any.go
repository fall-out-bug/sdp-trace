package harnessobs

func findByKeyInSlice[T any](value []any, wanted map[string]bool, match func(any) (T, bool)) (T, bool) {
	// findByKeyInSlice keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var zero T
	for _, child := range value {

		if found, ok := findByKeyIn(child, wanted, match); ok {
			return found, true
		}
	}
	return zero, false
}
