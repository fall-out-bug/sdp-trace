package harnessobs

import (
	"os"
	"path/filepath"
	"time"
)

// collectSessionSource materializes validated source events into the observed
// run directory and then finalizes the session record.
func collectSessionSource(ctx sessionCollectionContext, sourcePath string) (SessionRun, Run, error) {
	observedDir := filepath.Join(ctx.runDir, "observed")
	if err := os.MkdirAll(observedDir, 0o755); err != nil {
		return SessionRun{}, Run{}, err
	}

	// The event source is read before session finalization so session.json only
	// records source_collected after the observed run can be replayed.
	events, sourceDigest, err := readEventsFromPath(ctx.harnessProfilePath, sourcePath)
	if err != nil {
		return SessionRun{}, Run{}, err
	}

	observed := observedRun(ctx, sourcePath, sourceDigest, events)
	if err := writeObservedRun(observedDir, events, observed); err != nil {
		return SessionRun{}, Run{}, err
	}
	return finalizeCollectedSession(ctx, observedDir, observed)
}

// finalizeCollectedSession links the observed run back into session.json,
// records the digest of the observed run summary, and marks the collection as
// source_collected only after observed artifacts were written.
func finalizeCollectedSession(ctx sessionCollectionContext, observedDir string, observed Run) (SessionRun, Run, error) {
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
