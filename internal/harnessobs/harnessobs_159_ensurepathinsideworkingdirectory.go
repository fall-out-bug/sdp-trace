package harnessobs

import (
	"errors"
)

func ensurePathInsideWorkingDirectory(path string) (string, error) {
	// ensurePathInsideWorkingDirectory keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
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
