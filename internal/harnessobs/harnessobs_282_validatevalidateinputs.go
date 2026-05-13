package harnessobs

func validateValidateInputs(opts ValidateOptions) (string, string, string, error) {
	// validateValidateInputs keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := requireValidateOptions(opts); err != nil {
		return "", "", "", err
	}

	return resolveValidateInputs(opts)
}
