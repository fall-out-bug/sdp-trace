package harnessobs

import "fmt"

// resolveValidateInputs returns normalized safe paths for validation while
// preserving the optional output path contract.
func resolveValidateInputs(opts ValidateOptions) (string, string, string, error) {
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

// resolveValidateSourcePaths validates the profile file and run directory
// independently so callers get profile/run-specific unsafe path errors.
func resolveValidateSourcePaths(opts ValidateOptions) (string, string, error) {
	profilePath, err := safeExistingFile(opts.ProfilePath)
	if err != nil {
		return "", "", fmt.Errorf("unsafe profile path: %w", err)
	}
	runDir, err := safeExistingDir(opts.RunDir)
	if err != nil {
		return "", "", fmt.Errorf("unsafe run path: %w", err)
	}
	return profilePath, runDir, nil
}

// resolveValidateOutPath keeps --out optional, but applies output-file safety
// whenever the caller requests a validation artifact.
func resolveValidateOutPath(outPath string) (string, error) {
	if outPath == "" {
		return "", nil
	}
	safeOut, err := safeOutFile(outPath)
	if err != nil {
		return "", fmt.Errorf("unsafe out path: %w", err)
	}
	return safeOut, nil
}
