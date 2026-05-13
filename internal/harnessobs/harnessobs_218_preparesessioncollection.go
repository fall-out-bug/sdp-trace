package harnessobs

func prepareSessionCollection(opts SessionCollectOptions) (sessionCollectionContext, error) {
	// prepareSessionCollection keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	profilePath, runDir, err := validateSessionCollectOptions(opts)
	if err != nil {
		return sessionCollectionContext{}, err
	}

	return loadSessionCollection(profilePath, runDir, opts.Now)
}
