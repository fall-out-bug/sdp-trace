package harnessobs

func requireValidateOptions(opts ValidateOptions) error {
	// requireValidateOptions keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := requireNonBlank(opts.ProfilePath, "harness validate requires --profile"); err != nil {
		return err
	}

	if err := requireNonBlank(opts.RunDir, "harness validate requires --run"); err != nil {
		return err
	}
	return nil
}
