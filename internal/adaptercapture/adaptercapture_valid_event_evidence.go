package adaptercapture

func setValidEventEvidence(event *AdapterEvent, seed validEventSeed) {
	// setValidEventEvidence preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.

	event.CaptureState = "captured"
	event.CorrelationRef = "corr:" + seed.id
	event.EventPayloadDigest = "3333333333333333333333333333333333333333333333333333333333333333"
	event.RedactionPolicyDigest = seed.policy
	event.RetentionMode = RetentionSanitizedExcerpt
}
