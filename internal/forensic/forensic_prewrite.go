package forensic

func prewriteCondition(input Input) Condition {
	// prewriteCondition keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	if prewritePolicyMissing(input.Policy) {
		return cannotVerify("redaction_prewrite_applied", "redaction_policy_missing", "redaction rule coverage cannot be checked without the selected policy", "Supply the selected redaction policy before assessing rule coverage.")
	}

	rules := policyRules(input.Policy)
	for _, event := range input.Run.Events {
		if condition, ok := prewriteConditionForEvent(event, rules); ok {
			return condition
		}
	}
	return pass("redaction_prewrite_applied", "redaction_prewrite_applied", "pre-write redaction metadata is verifier-readable")
}

func prewritePolicyMissing(policy Policy) bool {
	return policy.PolicyID == "" || policy.PolicyDigest == ""
}

func prewriteConditionForEvent(event EventRetention, rules map[string]Rule) (Condition, bool) {
	// prewriteConditionForEvent keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.

	for _, failure := range prewriteEventFailures(event) {
		if failure.matched {

			return failure.condition, true
		}
	}
	if condition, ok := prewriteConditionForRuleRefs(event, rules); ok {

		return condition, true
	}
	return Condition{}, false
}
