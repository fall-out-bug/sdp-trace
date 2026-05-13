package authority

func validateTargetRules(env AuthorityEnvelope) string {
	// validateTargetRules keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	for _, rule := range env.TargetRules {

		if reason := validateTargetRule(env, rule); reason != "" {
			return reason
		}
	}
	return ""
}

func validateTargetRule(env AuthorityEnvelope, rule TargetRule) string {
	// validateTargetRule keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	if targetRuleMalformed(rule) {

		return "target_rule_malformed"
	}
	if validateEventSet(rule.AllowedEvents, rule.DeniedEvents) != "" {

		return "target_rule_conflict"
	}
	if targetRuleConflictsWithTopLevel(env, rule) {
		return "target_rule_conflicts_with_top_level"
	}
	return ""
}

func targetRuleMalformed(rule TargetRule) bool {
	return rule.RuleID == "" || rule.TargetPattern == ""
}

func targetRuleConflictsWithTopLevel(env AuthorityEnvelope, rule TargetRule) bool {
	return eventSetIntersects(rule.AllowedEvents, env.DeniedEvents) ||
		eventSetIntersects(rule.DeniedEvents, env.AllowedEvents)
}
