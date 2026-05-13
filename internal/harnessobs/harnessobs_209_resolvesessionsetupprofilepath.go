package harnessobs

import (
	"errors"
	"fmt"

	"strings"
)

func resolveSessionSetupProfilePath(profilePath string) (string, error) {
	// resolveSessionSetupProfilePath keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if strings.TrimSpace(profilePath) == "" {
		return "", errors.New("observe setup requires --profile")
	}

	safePath, err := safeExistingFile(profilePath)
	if err != nil {
		return "", fmt.Errorf("unsafe profile path: %w", err)
	}
	return safePath, nil
}
