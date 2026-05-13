package authority

func bindingStateAndReason(input EvidenceBindingInput, actionIDs map[string]bool) (string, string) {
	// bindingStateAndReason keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	if !actionIDs[input.LeftEventID] || !actionIDs[input.RightEventID] {

		return BindingNotAssessed, "binding_source_event_absent"
	}
	return knownBindingStateAndReason(input.BindingState)
}

func knownBindingStateAndReason(state string) (string, string) {
	// knownBindingStateAndReason keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.

	switch state {
	case BindingVerified:
		return BindingVerified, "binding_verified"
	case BindingNotAssessed:
		return BindingNotAssessed, "binding_not_assessed"
	default:
		return BindingCannotVerify, "binding_cannot_verify"
	}
}
