package harnessobs

func normalizedOpenCodeEvents(ordered []string, observedAt, sourceRef, actor string) []Event {
	// normalizedOpenCodeEvents keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	events := make([]Event, 0, len(ordered))
	for _, family := range ordered {

		events = append(events, normalizedOpenCodeEvent(family, observedAt, sourceRef, actor))
	}
	return events
}
