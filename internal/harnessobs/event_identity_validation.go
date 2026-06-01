package harnessobs

import (
	"errors"
	"time"
)

// Event identity validation keeps schema, family, type, and timestamp checks
// together because all guard whether an event belongs to the selected profile.
// Omitted timestamps remain valid, but present timestamps must be RFC3339 so
// evidence summaries do not imply local-time interpretation.
func validateEventIdentity(profile Profile, event Event) error {
	for _, check := range eventIdentityChecks(profile, event) {
		if !check.ok {
			return errors.New(check.err)
		}
	}
	return validateObservedAt(event.ObservedAt)
}

func eventIdentityChecks(profile Profile, event Event) []eventRefCheck {
	return []eventRefCheck{
		{safeFileIDPattern.MatchString(event.EventID), "unsafe event_id"},
		{event.EventSchemaVersion == profile.EventSchemaVersion, "schema_version_mismatch"},
		{validFamily(event.EventFamily), "unsupported event_family"},
		{safeIDPattern.MatchString(event.EventType), "unsafe event_type"},
	}
}

func validateObservedAt(value string) error {
	if !validObservedAt(value) {
		return errors.New("invalid observed_at")
	}
	return nil
}

func validObservedAt(value string) bool {
	if value == "" {
		return true
	}
	_, err := time.Parse(time.RFC3339, value)
	return err == nil
}
