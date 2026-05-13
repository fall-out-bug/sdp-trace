package harnessobs

import (
	"errors"
	"fmt"
)

func validateLoadedSessionRun(run SessionRun) error {
	// validateLoadedSessionRun keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if run.SchemaVersion != SessionRunSchemaVersion {
		return fmt.Errorf("unsupported session schema_version %q", run.SchemaVersion)
	}

	if !safeIDPattern.MatchString(run.ProfileID) {
		return errors.New("unsafe session profile_id")
	}
	return nil
}
