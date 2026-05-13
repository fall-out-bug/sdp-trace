package forensic

// ValidTestInput is used by CLI tests to write representative fixtures without
// duplicating Block 18 policy/run semantics outside this package.
func ValidTestInput() Input {
	return validTestInput()
}
func validTestInput() Input {
	// validTestInput keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	policyDigest := "1111111111111111111111111111111111111111111111111111111111111111"

	return Input{
		Policy: validTestPolicy(policyDigest),
		Run:    validRunEvidence(policyDigest),
	}
}

func validTestPolicy(policyDigest string) Policy {
	// validTestPolicy keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.

	return Policy{
		PolicyID:         "customer-forensic-policy-v1",
		SchemaVersion:    "1.0.0",
		PolicyDigest:     policyDigest,
		PolicyProvenance: Provenance{Source: "vcs", Digest: policyDigest},

		AllowedRetentionModes: validAllowedRetentionModes(),

		RedactionActions:            validRedactionActions(),
		ForbiddenPersistenceClasses: validForbiddenPersistenceClasses(),
		CriticalEventFamilies:       []string{"command_finished", "test_output_observed"},
		Authority:                   AuthorityRef{ActorID: "human:security-owner", VerificationState: AuthorityVerified},
		ProfileMappings:             validProfileMappings(),
		UnresolvedRedactionImpact:   "fail_forensic_retention",
		Rules:                       validRedactionRules(),
	}
}

func validAllowedRetentionModes() []string {
	return append([]string(nil), validAllowedRetentionModeFixture...)
}

func validRedactionActions() []string {
	return append([]string(nil), validRedactionActionFixture...)
}

func validForbiddenPersistenceClasses() []string {
	return append([]string(nil), validForbiddenPersistenceClassFixture...)
}
