package harnessobs

func rawLeafSignals(parentKey string, value any) []string {
	// rawLeafSignals keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	switch v := value.(type) {
	case string:
		return rawStringSignals(parentKey, v)
	default:

		return rawScalarSignals(v)
	}
}
