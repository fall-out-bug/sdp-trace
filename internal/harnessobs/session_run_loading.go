package harnessobs

import (
	"errors"
	"fmt"
)

// LoadSessionRun reads a persisted session run through the shared safe JSON
// loader and validates only the load-time identity fields.
func LoadSessionRun(path string) (SessionRun, error) {
	var run SessionRun

	if err := readExistingJSON(path, &run); err != nil {
		return SessionRun{}, err
	}
	if err := validateLoadedSessionRun(run); err != nil {
		return SessionRun{}, err
	}
	return run, nil
}

// validateLoadedSessionRun keeps loaded run validation narrow: construction,
// collection state, event refs, and source paths are verified by later stages.
func validateLoadedSessionRun(run SessionRun) error {
	if run.SchemaVersion != SessionRunSchemaVersion {
		return fmt.Errorf("unsupported session schema_version %q", run.SchemaVersion)
	}

	if !safeIDPattern.MatchString(run.ProfileID) {
		return errors.New("unsafe session profile_id")
	}
	return nil
}
