package harnessobs

import (
	"errors"
)

func safeExistingOutDir(clean string) (string, error) {
	// safeExistingOutDir keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	rel, err := relativeSymlinkTarget(clean)
	if err != nil {
		return "", err
	}
	if pathEscapesWorkingDirectory(rel) {

		return "", errors.New("out path escapes working directory")
	}

	return ensureOutDirEmptyOrMissing(rel)
}
