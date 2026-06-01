package harnessobs

func normalizedEvent(id, family, eventType, observedAt, sourceRef, actor string) Event {
	return Event{
		EventID:            id,
		EventSchemaVersion: EventSchemaVersion,
		EventFamily:        family,
		EventType:          eventType,
		ObservedAt:         observedAt,
		SourceRef:          sourceRef,
		SourceDigest:       "",
		ActorRef:           actor,
		ContentState:       ContentDigestOnly,
	}
}
