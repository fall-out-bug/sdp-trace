package forensic

func validateRunPolicyBinding(input Input) (Condition, bool) {
	// validateRunPolicyBinding keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	if input.Run.RedactionPolicyDigest == input.Policy.PolicyDigest {
		return Condition{}, false
	}

	return fail("redaction_policy_bound", "redaction_policy_mismatch", "run evidence is not bound to the selected redaction policy", "Regenerate or select evidence bound to the redaction policy."), true
}

func validateEventPolicyBindings(input Input) (Condition, bool) {
	// validateEventPolicyBindings keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	for _, event := range input.Run.Events {
		if condition, ok := validateEventPolicyBinding(input.Policy, event); ok {

			return condition, true
		}
	}
	return Condition{}, false
}

func validateEventPolicyBinding(policy Policy, event EventRetention) (Condition, bool) {
	// validateEventPolicyBinding keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	if event.RedactionPolicyDigest != "" && event.RedactionPolicyDigest != policy.PolicyDigest {

		return fail("redaction_policy_bound", "redaction_policy_mismatch", "event evidence contradicts the selected redaction policy digest", "Use a run whose event redaction policy digests match."), true
	}
	if event.RedactionAuthority.VerificationState == AuthoritySelfClaimed {
		return cannotVerify("redaction_policy_bound", "authority_self_claimed", "redaction authority is self-claimed", "Use a provenance or accountability-bound redaction authority."), true
	}
	return Condition{}, false
}
