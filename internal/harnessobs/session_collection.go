package harnessobs

import (
	"errors"
	"time"
)

// CollectSession preserves missing-source degradation as a collectable session
// state; every other source resolution error remains a hard failure.
func CollectSession(opts SessionCollectOptions) (SessionRun, Run, error) {
	ctx, err := prepareSessionCollection(opts)
	if err != nil {
		return SessionRun{}, Run{}, err
	}

	// Source resolution is the boundary between setup metadata and observed
	// evidence; only the explicit unavailable sentinel degrades to an artifact.
	sourcePath, err := resolveSessionEventSource(&ctx)
	if err != nil {
		if !errors.Is(err, errSessionSourceUnavailable) {
			return SessionRun{}, Run{}, err
		}
		// Missing sources still produce degraded session artifacts so downstream
		// validation can report cannot_verify instead of losing the run.
		return markSessionSourceUnavailable(ctx)
	}
	return collectSessionSource(ctx, sourcePath)
}

// prepareSessionCollection validates caller-controlled paths before loading
// profile and run artifacts from disk.
func prepareSessionCollection(opts SessionCollectOptions) (sessionCollectionContext, error) {
	profilePath, runDir, err := validateSessionCollectOptions(opts)
	if err != nil {
		return sessionCollectionContext{}, err
	}

	return loadSessionCollection(profilePath, runDir, opts.Now)
}

// loadSessionCollection binds the setup session to the harness profile before
// event-source resolution can add observed evidence.
func loadSessionCollection(profilePath, runDir string, now time.Time) (sessionCollectionContext, error) {
	profile, session, err := loadSessionCollectionInputs(profilePath, runDir)
	if err != nil {
		return sessionCollectionContext{}, err
	}

	// The harness profile can differ from the setup profile when collection
	// points at an external raw-event source profile.
	harnessProfilePath, harnessProfile, err := loadHarnessProfile(profilePath, profile)
	if err != nil {
		return sessionCollectionContext{}, err
	}
	return newSessionCollectionContext(
		profilePath,
		runDir,
		observationTime(now),
		profile,
		session,
		harnessProfilePath,
		harnessProfile,
	), nil
}
