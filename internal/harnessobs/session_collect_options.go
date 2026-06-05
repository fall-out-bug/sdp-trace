package harnessobs

import "fmt"

// validateSessionCollectOptions resolves the two user-supplied collection paths
// through the existing safe path gates before collection reads or writes data.
func validateSessionCollectOptions(opts SessionCollectOptions) (string, string, error) {
	if err := requireSessionCollectOptions(opts); err != nil {
		return "", "", err
	}

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

// requireSessionCollectOptions preserves the CLI-facing missing-flag messages
// before path validation adds filesystem-specific error context.
func requireSessionCollectOptions(opts SessionCollectOptions) error {
	if err := requireNonBlank(opts.ProfilePath, "observe collect requires --profile"); err != nil {
		return err
	}

	if err := requireNonBlank(opts.RunDir, "observe collect requires --run"); err != nil {
		return err
	}
	return nil
}
