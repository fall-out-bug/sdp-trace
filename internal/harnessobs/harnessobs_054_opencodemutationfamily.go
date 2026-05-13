package harnessobs

func openCodeMutationFamily(raw map[string]any, signals []string) bool {
	// openCodeMutationFamily keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return hasSignal(signals, "file.write", "file.edit", "file.patch", "file.delete", "mutation") ||
		hasSignalPrefix(signals, "mutation.") ||
		nativeMutationTool(raw)
}
