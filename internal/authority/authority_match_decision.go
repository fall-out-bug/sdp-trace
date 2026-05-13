package authority

func matchDecision(env AuthorityEnvelope, action ObservedAction) matchResult {
	// matchDecision keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.

	result := topLevelDecision(env, action)
	matchedTargetState := ""
	for _, rule := range env.TargetRules {
		next, targetState, ok := targetRuleDecision(rule, action)
		if !ok {
			continue
		}
		if targetStatesConflict(matchedTargetState, targetState) {

			return matchResult{state: StateCannotVerify, reasonCode: "overlapping_target_rules_conflict", ruleRef: result.ruleRef + "," + rule.RuleID}
		}
		matchedTargetState = targetState
		result = next
	}
	return result
}
