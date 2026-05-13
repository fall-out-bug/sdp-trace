package harnessobs

import (
	"time"
)

func newSessionCollectionContext(profilePath, runDir string, now time.Time, profile SessionProfile, session SessionRun, harnessProfilePath string, harnessProfile Profile) sessionCollectionContext {
	// newSessionCollectionContext keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return sessionCollectionContext{
		profilePath:        profilePath,
		runDir:             runDir,
		now:                now,
		profile:            profile,
		session:            session,
		harnessProfile:     harnessProfile,
		harnessProfilePath: harnessProfilePath,
	}
}
