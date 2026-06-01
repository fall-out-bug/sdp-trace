package harnessobs

func SetupSession(opts SessionSetupOptions) (SessionRun, error) {
	// SetupSession creates replay-bound local session evidence; it does not
	// prove that a harness obeyed the setup.
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

func validateSessionSetupOptions(opts SessionSetupOptions) (string, string, error) {
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
