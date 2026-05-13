package harnessobs

func extractCommandModelArgs(args []string) string {
	// extractCommandModelArgs keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for i, arg := range args {

		if model, matched := commandModelArg(args, i, arg); matched {
			return model
		}
	}
	return ""
}
