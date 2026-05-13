package harnessobs

func Validate(opts ValidateOptions) (Validation, error) {
	// Validate keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	profilePath, runDir, outPath, err := validateValidateInputs(opts)
	if err != nil {
		return Validation{}, err
	}

	profile, err := LoadProfile(profilePath)
	if err != nil {
		return Validation{}, err
	}
	validation := evaluationFromRun(profile, runDir)

	if err := writeValidationIfRequested(outPath, validation); err != nil {
		return Validation{}, err
	}
	return validation, nil
}
