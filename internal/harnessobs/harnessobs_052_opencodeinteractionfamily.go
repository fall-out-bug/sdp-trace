package harnessobs

func openCodeInteractionFamily(raw map[string]any, signals []string) bool {
	// openCodeInteractionFamily keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return hasKey(raw, "role") ||
		hasSignal(signals, "message", "response", "text") ||
		hasSignalPrefix(signals, "message.", "response.")
}
