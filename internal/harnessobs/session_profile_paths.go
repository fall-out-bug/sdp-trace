package harnessobs

import (
	"errors"
	"path/filepath"
	"strings"
)

// Profile-relative paths are declared inside session profiles and resolved
// against the profile file's directory before the existing input/output path
// policies inspect the concrete path.
func safeProfileRelativeFile(profilePath, relPath string) (string, error) {
	if unsafeProfileRelativePath(relPath) {
		return "", errors.New("profile relative path must be local without traversal")
	}

	baseDir := filepath.Dir(profilePath)
	if baseDir == "." {
		return safeExistingFile(relPath)
	}
	return safeExistingFile(filepath.Join(baseDir, relPath))
}

// Normalized event output still uses profile-relative addressing, but the
// final filesystem policy is output-file safety instead of existing-file
// safety because raw normalization materializes the file during collection.
func safeProfileRelativeOutFile(profilePath, relPath string) (string, error) {
	if unsafeProfileRelativePath(relPath) {
		return "", errors.New("profile relative output path must be local without traversal")
	}

	baseDir := filepath.Dir(profilePath)
	if baseDir == "." {
		return safeOutFile(relPath)
	}
	return safeOutFile(filepath.Join(baseDir, relPath))
}

// Profile declarations are intentionally stricter than generic output paths:
// they cannot be absolute, URL-like, or traversal-bearing before joining to
// the profile directory.
func unsafeProfileRelativePath(path string) bool {
	return filepath.IsAbs(path) || strings.Contains(path, "://") || strings.Contains(path, "..")
}
