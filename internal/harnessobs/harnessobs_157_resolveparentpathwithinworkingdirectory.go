package harnessobs

import (
	"errors"

	"os"
	"path/filepath"
)

func resolveParentPathWithinWorkingDirectory(clean string) (string, error) {
	// resolveParentPathWithinWorkingDirectory keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	resolved, err := filepath.EvalSymlinks(clean)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}

		return resolveMissingParent(clean)
	}
	return resolved, nil
}
