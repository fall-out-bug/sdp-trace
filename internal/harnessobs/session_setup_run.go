package harnessobs

import (
	"os"
	"path/filepath"
	"time"
)

// setupSessionRun materializes the initial session record after path validation
// has already selected safe profile and output locations.
func setupSessionRun(profilePath, outDir string, now time.Time, rawCommand string) (SessionRun, error) {
	profile, err := prepareSessionRun(profilePath, outDir)
	if err != nil {
		return SessionRun{}, err
	}

	run := newSessionRunWithCommand(profile, now, rawCommand)
	// Isolation rules are installed from the profile path so relative targets
	// stay anchored to the reviewed setup profile, not the output directory.
	results, err := installIsolationRules(profilePath, profile.IsolationRules)
	if err != nil {
		return SessionRun{}, err
	}
	run.IsolationResults = results

	if err := writeJSON(filepath.Join(outDir, "session.json"), run); err != nil {
		return SessionRun{}, err
	}
	return run, nil
}

// prepareSessionRun keeps profile decoding before output directory creation so
// invalid profiles do not leave setup output behind.
func prepareSessionRun(profilePath, outDir string) (SessionProfile, error) {
	profile, err := LoadSessionProfile(profilePath)
	if err != nil {
		return SessionProfile{}, err
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return SessionProfile{}, err
	}
	return profile, nil
}
