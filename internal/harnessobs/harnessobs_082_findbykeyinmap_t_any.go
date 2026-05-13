package harnessobs

func findByKeyInMap[T any](value map[string]any, wanted map[string]bool, match func(any) (T, bool)) (T, bool) {
	// findByKeyInMap keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var zero T
	for key, child := range value {

		if found, ok := matchWantedKey(key, child, wanted, match); ok {
			return found, true
		}
		if found, ok := findByKeyIn(child, wanted, match); ok {
			return found, true
		}
	}
	return zero, false
}
