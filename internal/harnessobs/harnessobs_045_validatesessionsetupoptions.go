package harnessobs

func validateSessionSetupOptions(opts SessionSetupOptions) (string, string, error) {
	// validateSessionSetupOptions keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	profilePath, err := resolveSessionSetupProfilePath(opts.ProfilePath)
	if err != nil {
		return "", "", err
	}
	outDir, err := resolveSessionSetupOutDir(opts.OutDir)
	if err != nil {
		return "", "", err
	}
	return profilePath, outDir, nil
}
