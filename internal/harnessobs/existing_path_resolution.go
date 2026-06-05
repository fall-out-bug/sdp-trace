package harnessobs

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func resolveExistingAbsolutePath(path string) (string, error) {
	// Symlinks are resolved before containment checks so a local-looking path
	// cannot point outside the working directory.
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}

	return absolutePath(resolved)
}

func sanitizeExistingPath(path, traversalError string) (string, error) {
	// Reject absolute, URL-like, and parent-traversal inputs before cleaning so
	// later filesystem calls only see local candidate paths.
	if filepath.IsAbs(path) || strings.Contains(path, "://") || strings.Contains(path, "..") {
		return "", errors.New(traversalError)
	}

	return filepath.Clean(path), nil
}

func relativeWorkingDirectoryPath(abs string) (string, error) {
	// The public validation returns working-directory-relative paths; escaping
	// paths are rejected instead of being represented with leading `..`.
	rel, err := relativePathFromWorkingDirectory(abs)
	if err != nil {
		return "", err
	}
	if pathEscapesWorkingDirectory(rel) {
		return "", errors.New("path escapes working directory")
	}
	return rel, nil
}

func ensureExpectedPathType(path string, spec existingPathSpec) error {
	// Type mismatches keep caller-specific messages from the entrypoint spec,
	// preserving file-vs-directory diagnostics.
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() == spec.requireDir {
		return nil
	}

	return errors.New(spec.typeError)
}
