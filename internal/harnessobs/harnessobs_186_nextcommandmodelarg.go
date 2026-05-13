package harnessobs

func nextCommandModelArg(args []string, i int) string {
	// nextCommandModelArg keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if i+1 >= len(args) {
		return ""
	}

	return safeCommandModel(args[i+1])
}
