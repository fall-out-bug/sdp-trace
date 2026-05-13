package harnessobs

func normalizedEvent(id, family, eventType, observedAt, sourceRef, actor string) Event {
	// normalizedEvent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

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
