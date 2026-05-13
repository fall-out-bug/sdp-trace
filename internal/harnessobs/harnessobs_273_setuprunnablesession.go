package harnessobs

func setupRunnableSession(opts SessionOptions) (SessionRun, error) {
	// setupRunnableSession keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := requireSessionCommand(opts.Command); err != nil {
		return SessionRun{}, err
	}

	return SetupSession(SessionSetupOptions{ProfilePath: opts.ProfilePath, OutDir: opts.OutDir, Now: opts.Now})
}
