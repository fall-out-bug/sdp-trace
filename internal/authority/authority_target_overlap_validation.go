package authority

func validateTargetRuleOverlap(rules []TargetRule) string {
	// validateTargetRuleOverlap keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	for i := range rules {

		if targetRuleConflictsWithAny(rules[i], rules[i+1:]) {
			return "overlapping_target_rules_conflict"
		}
	}
	return ""
}

func targetRuleConflictsWithAny(rule TargetRule, others []TargetRule) bool {
	// targetRuleConflictsWithAny keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	for _, other := range others {
		if targetRulesConflict(rule, other) {

			return true
		}
	}
	return false
}
