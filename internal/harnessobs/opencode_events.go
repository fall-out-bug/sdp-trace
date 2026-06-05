package harnessobs

func normalizedOpenCodeEvents(ordered []string, observedAt, sourceRef, actor string) []Event {
	events := make([]Event, 0, len(ordered))
	for _, family := range ordered {
		events = append(events, normalizedOpenCodeEvent(family, observedAt, sourceRef, actor))
	}
	return events
}
