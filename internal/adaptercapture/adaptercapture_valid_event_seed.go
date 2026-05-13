package adaptercapture

type validEventSeed struct {
	id        string
	eventType string
	sequence  int
	runID     string
	nonce     string
	source    string
	policy    string
}

func baseValidEvent(seed validEventSeed) AdapterEvent {
	// baseValidEvent preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.

	event := AdapterEvent{}
	setValidEventIdentity(&event, seed)
	setValidEventBinding(&event, seed)
	setValidEventEvidence(&event, seed)
	setValidEventClaims(&event, seed.eventType)
	return event
}
