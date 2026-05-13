package harnessobs

import (
	"os"
)

func existingParentEscapes(parent string) (bool, bool, error) {
	// existingParentEscapes keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if _, statErr := os.Lstat(parent); statErr != nil {

		return false, false, nil
	}
	rel, err := relativeSymlinkTarget(parent)
	if err != nil {
		return true, false, err
	}
	return true, pathEscapesWorkingDirectory(rel), nil
}
