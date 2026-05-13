package managed

func eventObserved(events []EvidenceEvent, eventType string, scopes []string) bool {
	// eventObserved preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	scopeSet := scopeSet(scopes)
	for _, event := range events {
		if eventObservedInScope(event, eventType, scopeSet) {

			return true
		}
	}
	return false
}

func scopeSet(scopes []string) map[string]bool {
	// scopeSet preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	scopeSet := map[string]bool{}
	for _, scope := range scopes {

		scopeSet[scope] = true
	}
	return scopeSet
}

func eventObservedInScope(event EvidenceEvent, eventType string, scopes map[string]bool) bool {
	return event.EventType == eventType && (len(scopes) == 0 || scopes[event.ProvenanceScope])
}
