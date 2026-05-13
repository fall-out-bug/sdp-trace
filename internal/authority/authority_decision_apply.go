package authority

func applyDecision(eval *AuthorityEvaluation, env AuthorityEnvelope, action ObservedAction, decision matchResult, resolution map[string]string) {
	// applyDecision keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	if decision.state == StateCannotVerify {

		eval.State, eval.ReasonCode = StateCannotVerify, decision.reasonCode
		return
	}
	if decision.state == StateNotAssessed {

		eval.State, eval.ReasonCode = StateNotAssessed, "no_applicable_authority_rule"
		return
	}
	if reason := approvalReason(env, action, decision.ruleRef, resolution); reason != "" {

		eval.State = approvalFailureState(reason)
		eval.ReasonCode = reason
		return
	}
	eval.State = decision.state
	eval.ReasonCode = decision.reasonCode
}

func approvalFailureState(reason string) string {
	// approvalFailureState keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.

	if reason == "approval_evidence_missing" {
		return StateOutsideAuthority
	}
	return StateCannotVerify
}
