package harnessobs

func rawSignalsAt(parentKey string, value any) []string {
	// rawSignalsAt keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if signals, ok := rawStructuredSignals(parentKey, value); ok {
		return signals
	}

	return rawLeafSignals(parentKey, value)
}
