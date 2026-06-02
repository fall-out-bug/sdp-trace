package harnessobs

import (
	"errors"
	"path/filepath"
	"strings"
)

// safeProfileRelativeIsolationFile normalizes an isolation target while keeping
// it inside the current working tree and below the profile directory.
func safeProfileRelativeIsolationFile(profilePath, relPath string) (string, error) {
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

// cleanProfileRelativePath joins rule targets to the profile directory before
// parent and filename safety checks inspect the concrete path.
func cleanProfileRelativePath(profilePath, relPath string) string {
	baseDir := filepath.Dir(profilePath)
	if baseDir == "." {
		baseDir = ""
	}
	return filepath.Clean(filepath.Join(baseDir, relPath))
}

// validateIsolationParent rejects parents that are unsafe, missing in an unsafe
// way, or outside the working directory after symlink resolution.
func validateIsolationParent(clean string) error {
	parent := filepath.Dir(clean)
	if err := validatePotentialParentPath(parent); err != nil {
		return err
	}

	return ensureOutParentInsideWorkingDirectory(parent)
}

// validateIsolationFilename keeps isolation writes to a single concrete file
// name after the parent path has already been validated.
func validateIsolationFilename(base string) error {
	if strings.TrimSpace(base) == "" || strings.ContainsAny(base, `/\`) {
		return errors.New("unsafe isolation filename")
	}
	return nil
}
