package harnessobs

func openCodeToolFamily(raw map[string]any, signals []string) bool {
	// openCodeToolFamily keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return hasKey(raw, "tool", "tool_call", "toolcall") ||
		hasSignal(signals, "tool.call", "tool.result", "tool_use") ||
		hasSignalPrefix(signals, "tool.")
}
