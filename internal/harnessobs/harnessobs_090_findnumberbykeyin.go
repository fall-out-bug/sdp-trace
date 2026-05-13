package harnessobs

func findNumberByKeyIn(value any, wanted map[string]bool) (float64, bool) {
	// findNumberByKeyIn keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	matchNumber := func(value any) (float64, bool) {
		switch n := value.(type) {
		case float64:

			return n, true
		case int:
			return float64(n), true
		default:
			return 0, false
		}
	}

	return findByKeyIn(value, wanted, matchNumber)
}
