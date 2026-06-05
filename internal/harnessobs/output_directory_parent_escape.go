package harnessobs

import (
	"errors"
	"os"
	"path/filepath"
)

// Missing output directories are allowed only when their nearest existing
// parent resolves inside the working directory.
func ensureOutParentInsideWorkingDirectory(clean string) error {
	parentEscapes, err := outParentEscapes(clean)
	if err != nil {
		return err
	}
	if parentEscapes {
		return errors.New("out parent path escapes working directory")
	}
	return nil
}

func outParentEscapes(clean string) (bool, error) {
	// Walk upward to the closest existing parent; missing descendants are safe
	// only if that existing ancestor remains inside the working directory.
	parent := filepath.Dir(clean)
	for parent != "." && parent != string(filepath.Separator) {
		found, escapes, err := existingParentEscapes(parent)
		if found {
			return escapes, err
		}
		parent = filepath.Dir(parent)
	}
	return false, nil
}

// existingParentEscapes reports whether parent exists and whether resolving
// that parent crosses the working-directory boundary through a symlink.
func existingParentEscapes(parent string) (bool, bool, error) {
	if _, statErr := os.Lstat(parent); statErr != nil {
		return false, false, nil
	}
	rel, err := relativeSymlinkTarget(parent)
	if err != nil {
		return true, false, err
	}
	return true, pathEscapesWorkingDirectory(rel), nil
}
