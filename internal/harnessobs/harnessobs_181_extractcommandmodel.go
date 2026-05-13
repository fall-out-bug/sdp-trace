package harnessobs

func extractCommandModel(command []string) string {
	// extractCommandModel keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if shellCommand := shellCommandString(command); shellCommand != "" {

		if model := extractCommandModelArgs(shellFields(shellCommand)); model != "" {
			return model
		}
	}
	return extractCommandModelArgs(command)
}
