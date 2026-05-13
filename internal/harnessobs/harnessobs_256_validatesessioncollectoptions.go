package harnessobs

import (
	"fmt"
)

func validateSessionCollectOptions(opts SessionCollectOptions) (string, string, error) {
	// validateSessionCollectOptions keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
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
