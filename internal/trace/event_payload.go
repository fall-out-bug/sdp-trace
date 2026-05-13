package trace

// WithComputedPayloadDigest computes PayloadDigest from the event payload bytes.
func (event Event) WithComputedPayloadDigest() (Event, error) {
	// Payload digest calculation first makes legacy EventPayload and raw Payload
	// representations agree.
	event = event.EnsureDefaults()
	synced, err := event.syncPayloadRepresentation()
	if err != nil {
		return Event{}, err
	}
	event = synced
	payloadDigest, err := CanonicalEventPayloadDigest(event.Payload)
	if err != nil {
		return Event{}, err
	}
	event.PayloadDigest = payloadDigest
	return event, nil
}

// WithComputedEventHash computes payload and event hashes.
func (event Event) WithComputedEventHash() (Event, error) {
	// The event hash is computed only after the payload digest is populated,
	// making the digest part of the event-chain integrity material.
	withDigest, err := event.WithComputedPayloadDigest()
	if err != nil {
		return Event{}, err
	}
	eventHash, err := ComputeEventHash(withDigest)
	if err != nil {
		return Event{}, err
	}
	withDigest.EventHash = eventHash
	return withDigest, nil
}
