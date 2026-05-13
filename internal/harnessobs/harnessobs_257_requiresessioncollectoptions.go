package harnessobs

func requireSessionCollectOptions(opts SessionCollectOptions) error {
	// requireSessionCollectOptions keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := requireNonBlank(opts.ProfilePath, "observe collect requires --profile"); err != nil {
		return err
	}

	if err := requireNonBlank(opts.RunDir, "observe collect requires --run"); err != nil {
		return err
	}
	return nil
}
