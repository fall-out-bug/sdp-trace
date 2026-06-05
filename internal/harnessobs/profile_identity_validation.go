package harnessobs

import (
	"errors"
	"fmt"
)

// Profile identity validation keeps schema and identifier checks together
// because both guard whether a profile can describe portable evidence.
func validateProfileIdentity(profile Profile) error {
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
