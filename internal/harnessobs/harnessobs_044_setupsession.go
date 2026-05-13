package harnessobs

func SetupSession(opts SessionSetupOptions) (SessionRun, error) {
	// SetupSession keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	profilePath, outDir, err := validateSessionSetupOptions(opts)
	if err != nil {
		return SessionRun{}, err
	}

	run, err := setupSessionRun(profilePath, outDir, opts.Now, opts.Command)
	if err != nil {
		return SessionRun{}, err
	}
	return run, nil
}
