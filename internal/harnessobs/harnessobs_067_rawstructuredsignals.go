package harnessobs

func rawStructuredSignals(parentKey string, value any) ([]string, bool) {
	// rawStructuredSignals keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	switch v := value.(type) {
	case map[string]any:

		return rawMapSignals(v), true
	case []any:

		return rawSliceSignals(parentKey, v), true
	default:
		return nil, false
	}
}
