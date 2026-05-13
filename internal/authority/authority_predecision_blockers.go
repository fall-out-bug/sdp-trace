package authority

func applyPreDecisionBlockers(eval *AuthorityEvaluation, env AuthorityEnvelope, envState, envReason string, action ObservedAction, eventBindings []EvidenceBinding, resolution map[string]string) bool {
	// applyPreDecisionBlockers keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.

	state, reason := preDecisionBlocker(env, envState, envReason, action, eventBindings, resolution)
	if reason == "" {
		return false
	}
	eval.State = state
	eval.ReasonCode = reason
	return true
}

func preDecisionBlocker(env AuthorityEnvelope, envState, envReason string, action ObservedAction, eventBindings []EvidenceBinding, resolution map[string]string) (string, string) {
	// preDecisionBlocker keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	if envState != "" {

		return envState, envReason
	}
	if !validEventType(action.EventType) {

		return StateCannotVerify, "unsupported_event_type"
	}
	if state, reason := preDecisionReason(env, action, eventBindings, resolution); reason != "" {
		return state, reason
	}
	return "", ""
}
