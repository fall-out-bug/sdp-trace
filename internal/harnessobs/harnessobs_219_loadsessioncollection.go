package harnessobs

import (
	"time"
)

func loadSessionCollection(profilePath, runDir string, now time.Time) (sessionCollectionContext, error) {
	// loadSessionCollection keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	profile, session, err := loadSessionCollectionInputs(profilePath, runDir)
	if err != nil {
		return sessionCollectionContext{}, err
	}

	harnessProfilePath, harnessProfile, err := loadHarnessProfile(profilePath, profile)
	if err != nil {
		return sessionCollectionContext{}, err
	}
	now = sessionCollectionTime(now)
	return newSessionCollectionContext(profilePath, runDir, now, profile, session, harnessProfilePath, harnessProfile), nil
}
