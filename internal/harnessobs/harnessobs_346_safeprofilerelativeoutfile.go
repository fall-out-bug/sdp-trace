package harnessobs

import (
	"errors"

	"path/filepath"
)

func safeProfileRelativeOutFile(profilePath, relPath string) (string, error) {
	// safeProfileRelativeOutFile keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if unsafeProfileRelativePath(relPath) {
		return "", errors.New("profile relative output path must be local without traversal")
	}

	baseDir := filepath.Dir(profilePath)
	if baseDir == "." {
		return safeOutFile(relPath)
	}
	return safeOutFile(filepath.Join(baseDir, relPath))
}
