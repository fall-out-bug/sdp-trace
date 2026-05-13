package harnessobs

import (
	"errors"
)

func validateProfileMetadata(profile Profile) error {
	// validateProfileMetadata keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := validateProfileIdentity(profile); err != nil {
		return err
	}
	if profile.EventSchemaVersion != EventSchemaVersion {
		return errors.New("unsupported event_schema_version")
	}

	if len(profile.RequiredEventFamilies) == 0 {
		return errors.New("profile requires at least one required_event_family")
	}
	return nil
}
