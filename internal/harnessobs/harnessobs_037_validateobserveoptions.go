package harnessobs

func validateObserveOptions(opts ObserveOptions) (string, string, string, error) {
	// validateObserveOptions keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := requireObserveOptions(opts); err != nil {
		return "", "", "", err
	}

	return resolveObservePaths(opts)
}
