package harnessobs

import (
	"os"
	"path/filepath"
)

func resolveMissingParent(clean string) (string, error) {
	// resolveMissingParent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	parent := filepath.Dir(clean)
	if parent == clean {

		return "", os.ErrNotExist
	}
	return resolveParentPathWithinWorkingDirectory(parent)
}
