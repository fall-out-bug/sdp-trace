package forensic

func prewriteConditionForRuleRefs(event EventRetention, rules map[string]Rule) (Condition, bool) {
	// prewriteConditionForRuleRefs keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	for _, ruleRef := range event.RedactionRuleRefs {

		if condition, ok := prewriteConditionForRuleRef(ruleRef, event.RedactionAction, rules); ok {
			return condition, true
		}
	}
	return Condition{}, false
}
func prewriteConditionForRuleRef(ruleRef, eventAction string, rules map[string]Rule) (Condition, bool) {
	// prewriteConditionForRuleRef keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	rule, ok := rules[ruleRef]
	if !ok {

		return fail("redaction_prewrite_applied", "redaction_rule_unknown", "event references a redaction rule that is absent from the selected redaction policy", "Use event redaction_rule_refs from the selected redaction policy."), true
	}
	if prewriteRuleActionMismatch(rule.Action, eventAction) {

		return fail("redaction_prewrite_applied", "redaction_rule_action_mismatch", "event redaction action contradicts the selected policy rule", "Align event redaction action with the selected policy rule."), true
	}
	return Condition{}, false
}

func prewriteRuleActionMismatch(ruleAction, eventAction string) bool {
	return ruleAction != "" && ruleAction != eventAction
}

func prewriteEventHasSecretLike(event EventRetention) bool {
	return event.SecretLikeValuePresent
}

func prewriteMissingRedactionDigests(event EventRetention) bool {
	return event.RedactionInputDigest == "" || event.RedactedPayloadDigest == ""
}

func prewriteRuleRefsMissing(event EventRetention) bool {
	return event.RedactionAction == RedactionActionApplyRule && len(event.RedactionRuleRefs) == 0
}
