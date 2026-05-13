package harnessobs

func appendScannedEvent(events []Event, event Event, ok bool) []Event {
	// appendScannedEvent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if ok {

		return append(events, event)
	}
	return events
}
