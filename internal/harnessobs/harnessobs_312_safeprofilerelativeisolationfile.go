package harnessobs

import (
	"errors"

	"path/filepath"
)

func safeProfileRelativeIsolationFile(profilePath, relPath string) (string, error) {
	// safeProfileRelativeIsolationFile keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if unsafeProfileRelativePath(relPath) {
		return "", errors.New("profile relative isolation path must be local without traversal")
	}

	clean := cleanProfileRelativePath(profilePath, relPath)
	if err := validateIsolationParent(clean); err != nil {
		return "", err
	}
	if err := validateIsolationFilename(filepath.Base(clean)); err != nil {
		return "", err
	}
	return clean, nil
}
