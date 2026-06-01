package harnessobs

import (
	"errors"
	"path/filepath"
	"strings"
)

func normalizePotentialParentPath(path string) (string, error) {
	// Empty parent paths mean the output lives in the current working
	// directory, but explicit unsafe parents must still be rejected.
	if path == "" {
		path = "."
	}
	if err := validatePotentialParentPath(path); err != nil {
		return "", err
	}
	return filepath.Clean(path), nil
}

func validatePotentialParentPath(path string) error {
	if unsafePotentialParentPath(path) {
		return errors.New("parent path must be relative local directory without traversal")
	}
	return nil
}

func unsafePotentialParentPath(path string) bool {
	return filepath.IsAbs(path) || strings.Contains(path, "://") || strings.Contains(path, "..")
}
