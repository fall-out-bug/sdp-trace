package harnessobs

import (
	"errors"
	"os"
	"path/filepath"
)

func ensurePathInsideWorkingDirectory(path string) (string, error) {
	// Containment is evaluated on an absolute path, then returned as a
	// working-directory-relative path for artifact writers.
	absPath, err := absolutePath(path)
	if err != nil {
		return "", err
	}
	rel, err := relativePathFromWorkingDirectory(absPath)
	if err != nil {
		return "", err
	}
	if pathEscapesWorkingDirectory(rel) {
		return "", errors.New("parent path escapes working directory")
	}

	return rel, nil
}

func absolutePath(path string) (string, error) {
	return filepath.Abs(path)
}

func relativePathFromWorkingDirectory(path string) (string, error) {
	// Use the live working directory because callers deliberately chdir in
	// tests and command execution before resolving local artifact paths.
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Rel(cwd, path)
}
