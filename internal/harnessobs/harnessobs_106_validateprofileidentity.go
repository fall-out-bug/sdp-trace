package harnessobs

import (
	"errors"
	"fmt"
)

func validateProfileIdentity(profile Profile) error {
	// validateProfileIdentity keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if profile.SchemaVersion != ProfileSchemaVersion {
		return fmt.Errorf("unsupported harness profile schema_version: %s", profile.SchemaVersion)
	}

	if !safeIDPattern.MatchString(profile.ProfileID) {
		return errors.New("unsafe profile_id")
	}

	if !safeIDPattern.MatchString(profile.HarnessFamily) {
		return errors.New("unsafe harness_family")
	}
	return nil
}
