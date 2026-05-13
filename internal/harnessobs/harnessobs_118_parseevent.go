package harnessobs

func parseEvent(profile Profile, line []byte, lineNo int) (Event, error) {
	// parseEvent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	event, err := decodeSafeEventLine(line, lineNo)
	if err != nil {
		return Event{}, err
	}

	return event, validateParsedEvent(profile, event, line, lineNo)
}
