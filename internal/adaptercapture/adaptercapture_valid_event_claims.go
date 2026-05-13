package adaptercapture

func setValidEventClaims(event *AdapterEvent, eventType string) {
	// setValidEventClaims preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.

	event.ActorAttributionState = "bound"
	event.ModelIdentityProvenance = "gateway_observed"
	event.TestProvenance = "ci_executed"
	event.ExecutedEvidenceClaimed = eventType == "test_observed"
	event.ToolFamily = "edit"
}
