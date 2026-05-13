package harnessobs

func resolveValidateInputs(opts ValidateOptions) (string, string, string, error) {
	// resolveValidateInputs keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	profilePath, runDir, err := resolveValidateSourcePaths(opts)
	if err != nil {
		return "", "", "", err
	}

	outPath, err := resolveValidateOutPath(opts.OutPath)
	if err != nil {
		return "", "", "", err
	}
	return profilePath, runDir, outPath, nil
}
