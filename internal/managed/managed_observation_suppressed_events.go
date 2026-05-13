package managed

func suppressedEventGroupCondition(input Input, id, group, reasonPrefix string) (Condition, bool) {
	// suppressedEventGroupCondition preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	suppressed, valid, satisfies := suppressionForGroup(input, group)
	if !suppressed {

		return Condition{}, false
	}
	return validSuppressedEventGroupCondition(id, group, reasonPrefix, valid, satisfies)
}

func validSuppressedEventGroupCondition(id, group, reasonPrefix string, valid, satisfies bool) (Condition, bool) {
	// validSuppressedEventGroupCondition preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	if valid && satisfies {

		return pass(id, reasonPrefix+"_event_suppressed_by_policy", "required "+group+" event is suppressed by policy for this profile"), true
	}
	if valid {

		return Condition{ID: id, State: StateSuppressed, ReasonCode: reasonPrefix + "_event_suppressed", Reason: "required " + group + " event is suppressed but does not satisfy the managed profile", NextAction: "Capture the required " + group + " event or authorize satisfying suppression in pre-run policy."}, true
	}
	return Condition{}, false
}
