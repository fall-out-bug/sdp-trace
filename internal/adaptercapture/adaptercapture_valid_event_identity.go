package adaptercapture

func setValidEventIdentity(event *AdapterEvent, seed validEventSeed) {
	// setValidEventIdentity preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.

	event.EventID = seed.id
	event.EventType = seed.eventType
	event.ProducerIdentity = "adapter:generic"
	event.AdapterIdentity = "adapter:generic"
	event.IdentityBinding = IdentityBound
	event.ProvenanceScope = "adapter_observed"
}
