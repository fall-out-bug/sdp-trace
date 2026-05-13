package authority

func preDecisionReason(env AuthorityEnvelope, action ObservedAction, eventBindings []EvidenceBinding, resolution map[string]string) (string, string) {
	// preDecisionReason keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	if state, reason := taskScopeReason(env, action); reason != "" {

		return state, reason
	}
	if reason := evidenceRefsReason(action.EvidenceRefs, resolution); reason != "" {
		return StateCannotVerify, reason
	}
	if bindingCannotVerify(eventBindings) {

		return StateCannotVerify, "evidence_binding_cannot_verify"
	}
	return "", ""
}

func taskScopeReason(env AuthorityEnvelope, action ObservedAction) (string, string) {
	// taskScopeReason keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	if env.AuthorityScope == "repository" {

		return "", ""
	}
	if action.TaskID == "" {

		return StateNotAssessed, "task_not_assessed"
	}
	if action.TaskID != env.TaskID {
		return StateNotAssessed, "task_outside_selected_envelope"
	}
	return "", ""
}
