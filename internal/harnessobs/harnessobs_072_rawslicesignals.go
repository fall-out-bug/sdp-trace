package harnessobs

func rawSliceSignals(parentKey string, values []any) []string {
	// rawSliceSignals keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	parts := make([]string, 0, len(values))
	for _, child := range values {

		parts = append(parts, rawSignalsAt(parentKey, child)...)
	}
	return parts
}
