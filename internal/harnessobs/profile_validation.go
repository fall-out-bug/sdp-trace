package harnessobs

import "errors"

// Profile validation coordinates metadata, event-family, and degradation-rule
// checks without changing the individual error contracts.
func validateProfile(profile Profile) error {
	if err := validateProfileMetadata(profile); err != nil {
		return err
	}
	if err := validateProfileEventFamilies(profile.RequiredEventFamilies, profile.OptionalEventFamilies); err != nil {
		return err
	}
	return validateProfileDegradationRules(profile.DegradationRules)
}

func validateProfileMetadata(profile Profile) error {
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
