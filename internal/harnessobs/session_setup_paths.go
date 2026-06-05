package harnessobs

import (
	"errors"
	"fmt"
	"strings"
)

// Setup paths are validated before the profile is loaded or output files are
// created, so unsafe input cannot leak into later session materialization.
func resolveSessionSetupProfilePath(profilePath string) (string, error) {
	if strings.TrimSpace(profilePath) == "" {
		return "", errors.New("observe setup requires --profile")
	}

	safePath, err := safeExistingFile(profilePath)
	if err != nil {
		return "", fmt.Errorf("unsafe profile path: %w", err)
	}
	return safePath, nil
}

// The setup output directory uses creation-path validation because it may not
// exist yet; existing-file validation would reject valid setup runs.
func resolveSessionSetupOutDir(outDir string) (string, error) {
	if strings.TrimSpace(outDir) == "" {
		return "", errors.New("observe setup requires --out")
	}

	return safeOutDir(outDir)
}
