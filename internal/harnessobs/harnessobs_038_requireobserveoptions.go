package harnessobs

func requireObserveOptions(opts ObserveOptions) error {
	// requireObserveOptions keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := requireNonBlank(opts.ProfilePath, "harness observe requires --profile"); err != nil {
		return err
	}

	if err := requireNonBlank(opts.SourcePath, "harness observe requires --source"); err != nil {
		return err
	}
	if err := requireNonBlank(opts.OutDir, "harness observe requires --out"); err != nil {
		return err
	}
	return nil
}
