package managed

func eventGroupCondition(input Input, id, group string) Condition {
	// Event group checks preserve missing, suppressed, and observed event states as
	// separate managed-mode evidence outcomes.
	required := eventTypesForGroup(input, group)
	if len(required) == 0 {

		return pass(id, "condition_pass", "no required events for group")
	}
	scopes := acceptableScopesForGroup(input, group)
	if !allEventsObserved(input.Run.ObservedEvents, required, scopes) {
		return missingEventGroupCondition(input, id, group)
	}
	return pass(id, "condition_pass", "required "+group+" events are observed")
}

func allEventsObserved(events []EvidenceEvent, required, scopes []string) bool {
	// allEventsObserved preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	for _, eventType := range required {
		if !eventObserved(events, eventType, scopes) {

			return false
		}
	}
	return true
}

func missingEventGroupCondition(input Input, id, group string) Condition {
	// missingEventGroupCondition preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	reasonPrefix := group
	if group == "file" {

		reasonPrefix = "file_mutation"
	}
	if condition, ok := suppressedEventGroupCondition(input, id, group, reasonPrefix); ok {
		return condition
	}
	return Condition{ID: id, State: StateMissingTelemetry, ReasonCode: reasonPrefix + "_event_missing", Reason: "required " + group + " event is missing", NextAction: "Run through a managed boundary that emits required " + group + " events."}
}
