package harnessobs

import (
	"errors"

	"path/filepath"
)

func loadSessionCollectionInputs(profilePath, runDir string) (SessionProfile, SessionRun, error) {
	// loadSessionCollectionInputs keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	profile, err := LoadSessionProfile(profilePath)
	if err != nil {
		return SessionProfile{}, SessionRun{}, err
	}

	session, err := LoadSessionRun(filepath.Join(runDir, "session.json"))
	if err != nil {
		return SessionProfile{}, SessionRun{}, err
	}
	if session.ProfileID != profile.ProfileID {

		return SessionProfile{}, SessionRun{}, errors.New("session profile mismatch")
	}

	return profile, session, nil
}
