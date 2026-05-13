package adaptercapture

func validEvent(id, eventType string, sequence int, runID, nonce, source, policy string) AdapterEvent {
	// validEvent preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.

	seed := validEventSeed{id: id, eventType: eventType, sequence: sequence, runID: runID, nonce: nonce, source: source, policy: policy}
	event := baseValidEvent(seed)
	if digestOnlyValidEvent(eventType) {

		event.RetentionMode = RetentionDigestOnly
	}
	return event
}
