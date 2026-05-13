package forensic

func policyCondition(input Input) Condition {
	// policyCondition keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.

	if condition, ok := validatePolicyContract(input.Policy); ok {
		return condition
	}
	if condition, ok := validateRunPolicyBinding(input); ok {
		return condition
	}
	if condition, ok := validateEventPolicyBindings(input); ok {
		return condition
	}
	return pass("redaction_policy_bound", "redaction_policy_bound", "redaction policy digest and authority evidence are bound")
}

func validatePolicyContract(policy Policy) (Condition, bool) {
	// validatePolicyContract keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	for _, check := range policyContractChecks(policy) {
		if check.failed {

			return check.condition, true
		}
	}
	return Condition{}, false
}
