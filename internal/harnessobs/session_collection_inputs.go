package harnessobs

import (
	"errors"
	"path/filepath"
)

// loadSessionCollectionInputs rejects profile/session mismatches before a
// collect run can attach new observed evidence to the wrong setup session.
func loadSessionCollectionInputs(profilePath, runDir string) (SessionProfile, SessionRun, error) {
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
