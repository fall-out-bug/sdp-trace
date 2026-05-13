package harnessobs

import (
	"errors"

	"path/filepath"
)

func safeProfileRelativeFile(profilePath, relPath string) (string, error) {
	// safeProfileRelativeFile keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if unsafeProfileRelativePath(relPath) {
		return "", errors.New("profile relative path must be local without traversal")
	}

	baseDir := filepath.Dir(profilePath)
	if baseDir == "." {
		return safeExistingFile(relPath)
	}
	return safeExistingFile(filepath.Join(baseDir, relPath))
}
