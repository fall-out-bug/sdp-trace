package harnessobs

import (
	"path/filepath"

	"time"
)

func setupSessionRun(profilePath, outDir string, now time.Time, rawCommand string) (SessionRun, error) {
	// setupSessionRun keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	profile, err := prepareSessionRun(profilePath, outDir)
	if err != nil {
		return SessionRun{}, err
	}

	run := newSessionRunWithCommand(profile, now, rawCommand)
	results, err := installIsolationRules(profilePath, profile.IsolationRules)
	if err != nil {
		return SessionRun{}, err
	}
	run.IsolationResults = results

	if err := writeSessionJSON(filepath.Join(outDir, "session.json"), run); err != nil {
		return SessionRun{}, err
	}
	return run, nil
}
