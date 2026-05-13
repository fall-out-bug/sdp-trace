package harnessobs

import (
	"errors"
)

func relativeWorkingDirectoryPath(abs string) (string, error) {
	// relativeWorkingDirectoryPath keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	rel, err := relativePathFromWorkingDirectory(abs)
	if err != nil {
		return "", err
	}
	if pathEscapesWorkingDirectory(rel) {

		return "", errors.New("path escapes working directory")
	}
	return rel, nil
}
