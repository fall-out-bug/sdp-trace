package harnessobs

import (
	"errors"
	"os"
	"path/filepath"
)

func safeOutDir(path string) (string, error) {
	// Output directories are local artifact roots; they may be created later,
	// but traversal and URL-like paths are rejected before filesystem checks.
	if unsafeOutputPath(path) {
		return "", errors.New("out must be a relative local directory without traversal")
	}

	clean := filepath.Clean(path)
	return safeCleanOutDir(clean)
}

func safeCleanOutDir(clean string) (string, error) {
	exists, err := pathExistsForLstat(clean)
	if err != nil {
		return "", err
	}
	if exists {
		return safeExistingOutDir(clean)
	}

	if err := ensureOutParentInsideWorkingDirectory(clean); err != nil {
		return "", err
	}
	return ensureOutDirEmptyOrMissing(clean)
}

func pathExistsForLstat(path string) (bool, error) {
	// Lstat detects existing symlink outputs without following them; the
	// symlink target is resolved only after this existence branch is selected.
	if _, err := os.Lstat(path); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}
