package harnessobs

// Validate resolves trusted local inputs, evaluates the observed run, and only
// writes a validation artifact when the caller requested an output path.
func Validate(opts ValidateOptions) (Validation, error) {
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

// writeValidationIfRequested keeps validation previews side-effect free unless
// an explicit safe output path was supplied.
func writeValidationIfRequested(outPath string, validation Validation) error {
	if outPath == "" {
		return nil
	}
	return writeJSON(outPath, validation)
}
