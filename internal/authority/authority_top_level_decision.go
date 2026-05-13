package authority

func topLevelDecision(env AuthorityEnvelope, action ObservedAction) matchResult {
	// topLevelDecision keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	if contains(env.DeniedEvents, action.EventType) {

		return matchResult{state: StateOutsideAuthority, reasonCode: "event_denied", ruleRef: "denied_events"}
	}
	if contains(env.AllowedEvents, action.EventType) {

		return matchResult{state: StateWithinAuthority, reasonCode: "event_allowed", ruleRef: "allowed_events"}
	}
	return matchResult{state: StateNotAssessed, reasonCode: "no_applicable_authority_rule"}
}

func targetRuleDecision(rule TargetRule, action ObservedAction) (matchResult, string, bool) {
	// targetRuleDecision keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	if !targetMatches(rule.TargetPattern, action.Target) {

		return matchResult{}, "", false
	}
	if contains(rule.DeniedEvents, action.EventType) {

		return matchResult{state: StateOutsideAuthority, reasonCode: "target_event_denied", ruleRef: rule.RuleID}, StateOutsideAuthority, true
	}
	if contains(rule.AllowedEvents, action.EventType) {

		return matchResult{state: StateWithinAuthority, reasonCode: "target_event_allowed", ruleRef: rule.RuleID}, StateWithinAuthority, true
	}
	return matchResult{}, "", false
}

func targetStatesConflict(previous, next string) bool {
	return (previous == StateOutsideAuthority && next == StateWithinAuthority) ||
		(previous == StateWithinAuthority && next == StateOutsideAuthority)
}
