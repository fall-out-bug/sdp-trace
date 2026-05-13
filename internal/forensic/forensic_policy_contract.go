package forensic

func policyContractChecks(policy Policy) []policyContractCheck {
	// policyContractChecks keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.

	return []policyContractCheck{
		{

			failed:    policy.PolicyID == "" || policy.PolicyDigest == "",
			condition: cannotVerify("redaction_policy_bound", "missing_redaction_policy", "redaction policy is required", "Supply a redaction policy with stable id, version, digest, and provenance."),
		},
		{

			failed:    policyContractIncomplete(policy),
			condition: cannotVerify("redaction_policy_bound", "redaction_policy_contract_incomplete", "redaction policy contract is incomplete", "Supply redaction actions, forbidden persistence classes, authority, profile mappings, and unresolved-redaction impact."),
		},
		{

			failed:    policy.Authority.VerificationState == AuthoritySelfClaimed,
			condition: cannotVerify("redaction_policy_bound", "authority_self_claimed", "redaction policy authority is self-claimed", "Use a provenance or accountability-bound redaction policy authority."),
		},
	}
}

func policyContractIncomplete(policy Policy) bool {
	// policyContractIncomplete keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.

	missingParts := []bool{
		len(policy.RedactionActions) == 0,
		len(policy.ForbiddenPersistenceClasses) == 0,
		policy.Authority.ActorID == "",
		len(policy.ProfileMappings) == 0,
		policy.UnresolvedRedactionImpact == "",
	}
	for _, missing := range missingParts {
		if missing {

			return true
		}
	}
	return false
}
