package harnessobs

import "time"

// newSessionCollectionContext carries both setup profile identity and harness
// profile identity so later source collection can report unavailable evidence
// without losing provenance.
func newSessionCollectionContext(profilePath, runDir string, now time.Time, profile SessionProfile, session SessionRun, harnessProfilePath string, harnessProfile Profile) sessionCollectionContext {
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
