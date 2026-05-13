package harnessobs

func safeEvent(eventID string) string {
	// safeEvent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if safeIDPattern.MatchString(eventID) {
		return eventID
	}

	return "unknown_event"
}
