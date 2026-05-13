package harnessobs

import (
	"os"
)

func prepareSessionRun(profilePath, outDir string) (SessionProfile, error) {
	// prepareSessionRun keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	profile, err := LoadSessionProfile(profilePath)
	if err != nil {
		return SessionProfile{}, err
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return SessionProfile{}, err
	}
	return profile, nil
}
