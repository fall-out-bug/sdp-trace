package harnessobs

import (
	"errors"
	"os"
	"path/filepath"
)

func resolveParentPathWithinWorkingDirectory(clean string) (string, error) {
	// Missing output parents resolve through their nearest existing ancestor;
	// existing parents are symlink-resolved before containment is checked.
	resolved, err := filepath.EvalSymlinks(clean)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}

		return resolveMissingParent(clean)
	}
	return resolved, nil
}

func resolveMissingParent(clean string) (string, error) {
	parent := filepath.Dir(clean)
	if parent == clean {
		return "", os.ErrNotExist
	}
	return resolveParentPathWithinWorkingDirectory(parent)
}
