package harnessobs

import (
	"path/filepath"

	"time"
)

func finalizeCollectedSession(ctx sessionCollectionContext, observedDir string, observed Run) (SessionRun, Run, error) {
	// finalizeCollectedSession keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	session := ctx.session
	session.ObservedRunDir = filepath.ToSlash("observed")
	session.OutputDigest = digestFile(filepath.Join(observedDir, "run.json"))

	session.CollectionState = StatePass
	session.CollectionReason = "source_collected"

	if session.EndTime == "" {
		session.EndTime = ctx.now.Format(time.RFC3339)
	}
	if err := writeJSON(filepath.Join(ctx.runDir, "session.json"), session); err != nil {
		return SessionRun{}, Run{}, err
	}
	return session, observed, nil
}
