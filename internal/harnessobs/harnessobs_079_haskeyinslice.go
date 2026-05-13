package harnessobs

func hasKeyInSlice(values []any, wanted map[string]bool) bool {
	// hasKeyInSlice keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for _, child := range values {

		if hasKeyIn(child, wanted) {
			return true
		}
	}
	return false
}
