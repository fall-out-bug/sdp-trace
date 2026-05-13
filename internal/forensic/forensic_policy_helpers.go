package forensic

func validRetentionMode(mode string) bool {
	return validRetentionModes[mode]
}
func policyRules(policy Policy) map[string]Rule {
	// policyRules keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	out := map[string]Rule{}
	for _, rule := range policy.Rules {
		if rule.RuleID != "" {

			out[rule.RuleID] = rule
		}
	}
	return out
}

func allowedRetentionModes(policy Policy) map[string]bool {
	// allowedRetentionModes keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	out := map[string]bool{}
	for _, mode := range policy.AllowedRetentionModes {

		out[mode] = true
	}
	return out
}
