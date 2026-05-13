package harnessobs

import (
	"errors"
	"fmt"
)

func validateSessionProfileIdentity(profile SessionProfile) error {
	// validateSessionProfileIdentity keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if profile.SchemaVersion != SessionProfileSchemaVersion {
		return fmt.Errorf("unsupported session profile schema_version %q", profile.SchemaVersion)
	}

	if !safeIDPattern.MatchString(profile.ProfileID) {
		return errors.New("unsafe session profile_id")
	}
	return validateSessionProfilePaths(profile)
}
