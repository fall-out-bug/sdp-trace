package harnessobs

import (
	"fmt"
)

func loadHarnessProfile(profilePath string, profile SessionProfile) (string, Profile, error) {
	// loadHarnessProfile keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	harnessProfilePath, err := safeProfileRelativeFile(profilePath, profile.HarnessProfilePath)
	if err != nil {
		return "", Profile{}, fmt.Errorf("unsafe harness profile path: %w", err)
	}
	harnessProfile, err := LoadProfile(harnessProfilePath)
	if err != nil {
		return "", Profile{}, err
	}
	return harnessProfilePath, harnessProfile, nil
}
