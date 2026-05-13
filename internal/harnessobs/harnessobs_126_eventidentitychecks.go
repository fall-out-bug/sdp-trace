package harnessobs

func eventIdentityChecks(profile Profile, event Event) []eventRefCheck {
	// eventIdentityChecks keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return []eventRefCheck{
		{safeFileIDPattern.MatchString(event.EventID), "unsafe event_id"},
		{event.EventSchemaVersion == profile.EventSchemaVersion, "schema_version_mismatch"},
		{validFamily(event.EventFamily), "unsupported event_family"},
		{safeIDPattern.MatchString(event.EventType), "unsafe event_type"},
	}
}
