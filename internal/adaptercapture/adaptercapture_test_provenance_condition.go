package adaptercapture

func testProvenanceCondition(run RunEvidence) Condition {
	// Test provenance evidence is evaluated separately because test mentions are not
	// proof that tests actually executed.
	if event, ok := firstEvent(run.AdapterEvents, "test_observed"); ok {

		return testProvenanceEventCondition(event)
	}
	if hasRequired(run, "test_observed") {
		return Condition{ID: "test_provenance_not_overclaimed", State: StateMissingTelemetry, ReasonCode: "test_event_missing", Reason: "required test adapter event is missing", NextAction: "Capture test_observed adapter evidence."}
	}
	return pass("test_provenance_not_overclaimed", "test_provenance_not_required", "test provenance was not required")
}

func firstEvent(events []AdapterEvent, eventType string) (AdapterEvent, bool) {
	// firstEvent preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.
	for _, event := range events {
		if event.EventType == eventType {

			return event, true
		}
	}
	return AdapterEvent{}, false
}

func testProvenanceEventCondition(event AdapterEvent) Condition {
	// testProvenanceEventCondition preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.
	if testProvenanceExecuted(event.TestProvenance) {
		return pass("test_provenance_not_overclaimed", "test_provenance_executed", "test evidence is bound to CI or wrapper execution")
	}

	return nonExecutedTestProvenanceCondition(event)
}
