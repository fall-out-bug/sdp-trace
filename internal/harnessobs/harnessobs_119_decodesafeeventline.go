package harnessobs

func decodeSafeEventLine(line []byte, lineNo int) (Event, error) {
	// decodeSafeEventLine keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	raw, err := decodeRawEventLine(line, lineNo)
	if err != nil {
		return Event{}, err
	}

	if err := rejectUnsafeEvent(raw, lineNo); err != nil {
		return Event{}, err
	}
	return decodeEventLine(line, lineNo)
}
