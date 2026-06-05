package harnessobs

import (
	"errors"
	"fmt"
)

func validateSessionProfile(profile *SessionProfile) error {
	// Keep validation ordered so identity/path failures are reported before
	// mutating defaults such as stream capture normalization.
	// Later setup and isolation validators own their detailed rule semantics.
	if err := validateSessionProfileIdentity(*profile); err != nil {
		return err
	}

	if err := normalizeSessionStreamCapture(profile); err != nil {
		return err
	}
	if err := validateSessionSetupActions(profile.SetupActions); err != nil {
		return err
	}
	return validateSessionIsolationRules(profile.IsolationRules)
}

func validateSessionProfileIdentity(profile SessionProfile) error {
	// Schema and ID checks guard the profile identity before any path-like
	// fields or raw event declarations are interpreted.
	if profile.SchemaVersion != SessionProfileSchemaVersion {
		return fmt.Errorf("unsupported session profile schema_version %q", profile.SchemaVersion)
	}

	if !safeIDPattern.MatchString(profile.ProfileID) {
		return errors.New("unsafe session profile_id")
	}
	return validateSessionProfilePaths(profile)
}

func validateSessionProfilePaths(profile SessionProfile) error {
	// Required normalized event paths are checked before optional raw-event
	// configuration so missing core evidence wins the error order.
	if err := validateRequiredSessionPaths(profile); err != nil {
		return err
	}

	return validateRawEventConfig(profile)
}

func validateRequiredSessionPaths(profile SessionProfile) error {
	// The shared non-blank helper preserves whitespace-only rejection across
	// validation and session setup option gates.
	if err := requireNonBlank(profile.HarnessProfilePath, "session profile requires harness_profile_path"); err != nil {
		return err
	}

	if err := requireNonBlank(profile.EventSourcePath, "session profile requires event_source_path"); err != nil {
		return err
	}
	return nil
}
