package harnessobs

import (
	"os"
	"path/filepath"
)

func collectSessionSource(ctx sessionCollectionContext, sourcePath string) (SessionRun, Run, error) {
	// collectSessionSource keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	observedDir := filepath.Join(ctx.runDir, "observed")
	if err := os.MkdirAll(observedDir, 0o755); err != nil {
		return SessionRun{}, Run{}, err
	}

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
