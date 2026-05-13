package managed

func testProvenanceCondition(input Input) Condition {
	// Test provenance evidence is evaluated independently because observed events
	// do not prove tests actually ran.
	if eventObserved(input.Run.TestEvidence, "test_observed", []string{"local_observed", "ci_witnessed"}) {

		return pass("test_provenance_not_agent_reported", "test_provenance_not_agent_reported", "test evidence is wrapper or CI observed")
	}
	if eventObserved(input.Run.TestEvidence, "test_observed", []string{"agent_reported"}) {

		return fail("test_provenance_not_agent_reported", "agent_reported_test_not_executed", "test evidence is only agent-reported", "Record test execution through the managed wrapper or CI.")
	}
	return Condition{ID: "test_provenance_not_agent_reported", State: StateMissingTelemetry, ReasonCode: "test_provenance_missing", Reason: "test execution evidence is missing", NextAction: "Record test execution through the managed wrapper or CI."}
}

func suppressionCondition(input Input) Condition {
	// suppressionCondition preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	for _, suppressed := range input.Run.SuppressedEventGroups {
		if !suppressionVerified(input.Policy, suppressed) {

			return fail("suppression_policy_valid", "suppression_unverified", "suppression is not authorized by pre-run policy provenance", "Use pre-run policy authority for suppression or capture the event.")
		}
	}
	return pass("suppression_policy_valid", "suppression_policy_valid", "suppression policy is valid or no suppression is present")
}

func suppressionVerified(policy Policy, suppressed SuppressedEventGroup) bool {
	// suppressionVerified preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	rule, ok := suppressionRuleForGroup(policy, suppressed.EventGroup)

	return ok &&
		suppressed.AuthorizedByPolicy &&
		preRunProvenance(suppressed.PolicyProvenanceSource) &&
		preRunProvenance(rule.PolicyProvenanceSource)
}
func bypassCondition(input Input) Condition {
	// Bypass evidence remains explicit so an intentional bypass cannot look like
	// a passing managed-adapter observation.
	if input.Run.BypassObserved {
		return fail("bypass_not_observed", "bypass_observed", "managed boundary bypass was observed", "Rerun without bypass or lower the claim.")
	}

	return pass("bypass_not_observed", "bypass_not_observed", "no managed boundary bypass is observed")
}
