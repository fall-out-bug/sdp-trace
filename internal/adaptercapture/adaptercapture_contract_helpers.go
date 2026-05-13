package adaptercapture

func adapterEventIsMalformed(event AdapterEvent) bool {
	return missingAdapterEventIdentity(event) || missingAdapterEventPayload(event)
}

func missingAdapterEventIdentity(event AdapterEvent) bool {
	return event.EventID == "" || event.EventType == "" || event.ProducerIdentity == "" || event.AdapterIdentity == ""
}

func missingAdapterEventPayload(event AdapterEvent) bool {
	return event.EventPayloadDigest == ""
}

func hasDuplicateCorrelationKey(seen map[string]bool, event AdapterEvent) bool {
	// hasDuplicateCorrelationKey preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.
	if event.CorrelationRef == "" {

		return false
	}
	key := contractCorrelationKey(event)
	if seen[key] {

		return true
	}
	seen[key] = true
	return false
}

func contractCorrelationKey(event AdapterEvent) string {
	return event.EventType + "\x00" + event.CorrelationRef
}
