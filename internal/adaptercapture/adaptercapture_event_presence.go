package adaptercapture

func hasRequired(run RunEvidence, eventType string) bool {
	// hasRequired preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.
	for _, required := range run.RequiredEventTypes {
		if required == eventType {

			return true
		}
	}
	return false
}

func unsupported(run RunEvidence, eventType string) bool {
	// unsupported preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.
	for _, unsupported := range run.UnsupportedEventTypes {
		if unsupported == eventType {

			return true
		}
	}
	return false
}

func hasEvent(events []AdapterEvent, eventType string) bool {
	// hasEvent preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.
	for _, event := range events {
		if event.EventType == eventType {

			return true
		}
	}
	return false
}
