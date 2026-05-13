package harnessobs

import (
	"path/filepath"
)

func markSessionSourceUnavailable(ctx sessionCollectionContext) (SessionRun, Run, error) {
	// markSessionSourceUnavailable keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	session := unavailableSession(ctx)

	if err := writeJSON(filepath.Join(ctx.runDir, "session.json"), session); err != nil {
		return SessionRun{}, Run{}, err
	}

	return session, unavailableObservedRun(ctx), nil
}
