package authority

func evaluateAction(evaluationID, selectedPolicyID string, env AuthorityEnvelope, envState, envReason string, action ObservedAction, eventBindings []EvidenceBinding, resolution map[string]string) AuthorityEvaluation {
	// evaluateAction keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.

	eval := newAuthorityEvaluation(evaluationID, selectedPolicyID, action, eventBindings)
	eval.MissingAttributes = missingAttributes(eval)
	if applyPreDecisionBlockers(&eval, env, envState, envReason, action, eventBindings, resolution) {
		return eval
	}
	decision := matchDecision(env, action)
	eval.MatchedRuleRef = decision.ruleRef
	applyDecision(&eval, env, action, decision, resolution)
	return eval
}

func newAuthorityEvaluation(evaluationID, selectedPolicyID string, action ObservedAction, eventBindings []EvidenceBinding) AuthorityEvaluation {
	// newAuthorityEvaluation keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.

	eval := baseAuthorityEvaluation(evaluationID, selectedPolicyID, action)
	applyBindingAttribution(&eval, action, eventBindings)
	return eval
}

func baseAuthorityEvaluation(evaluationID, selectedPolicyID string, action ObservedAction) AuthorityEvaluation {
	// baseAuthorityEvaluation keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.

	return AuthorityEvaluation{
		EvaluationID:     evaluationID,
		EventID:          action.EventID,
		PolicyID:         selectedPolicyID,
		ActorAttribution: actorAttributionState(action),
		ToolAttribution:  AttributionNotAssessed,
		ModelAttribution: AttributionNotAssessed,
		SourceCoverage:   uniqueStrings([]string{action.SourceType}),
		EvidenceRefs:     safeRefs(action.EvidenceRefs),
		ActorID:          action.ActorID,
		OperationID:      action.OperationID,
	}
}

func applyBindingAttribution(eval *AuthorityEvaluation, action ObservedAction, eventBindings []EvidenceBinding) {
	// applyBindingAttribution keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	if action.SourceType == "harness_log" && action.OperationID != "" {

		eval.ToolAttribution = AttributionVerified
	}
	if hasVerifiedGatewayBinding(action, eventBindings) {

		eval.ModelAttribution = AttributionVerified
	}
}
