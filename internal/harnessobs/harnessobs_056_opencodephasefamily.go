package harnessobs

func openCodePhaseFamily(raw map[string]any, signals []string) bool {
	// openCodePhaseFamily keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return hasKey(raw, "phase") ||
		hasSignal(signals, "phase") ||
		hasSignalPrefix(signals, "phase.", "gsd.", "gsd_")
}
