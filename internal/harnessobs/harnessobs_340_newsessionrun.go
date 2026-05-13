package harnessobs

import (
	"time"
)

func newSessionRun(profile SessionProfile, now time.Time) SessionRun {
	// newSessionRun keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	actionIDs := sessionSetupActionIDs(profile)
	commit, commitState := currentSourceCommitState()
	return newSessionRunRecord(profile, now, actionIDs, commit, commitState)
}
