package forensic

func validExternalArtifactEvent(policyDigest string) EventRetention {
	// validExternalArtifactEvent keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.

	return EventRetention{
		EventType:             "test_output_observed",
		RetentionMode:         RetentionModeExternalArtifactRef,
		ForensicImportance:    "critical",
		RedactionPolicyDigest: policyDigest,

		RedactionInputDigest:   "5555555555555555555555555555555555555555555555555555555555555555",
		RedactedPayloadDigest:  "6666666666666666666666666666666666666666666666666666666666666666",
		RedactionAction:        RedactionActionApplyRule,
		RedactionRuleRefs:      []string{"secret-token-v1"},
		SecretLikeValuePresent: false,
		RedactionAuthority:     AuthorityRef{ActorID: "human:security-owner", VerificationState: AuthorityVerified},

		RawReference: validExternalRawReference(),
	}
}
func validExternalRawReference() *RawReference {
	// validExternalRawReference keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.

	return &RawReference{
		ReferenceType:           RetentionModeExternalArtifactRef,
		ReferenceURI:            "artifact://ci/run-1/test-output",
		Digest:                  Digest{Algorithm: "sha256", Value: "7777777777777777777777777777777777777777777777777777777777777777"},
		AccessState:             AccessStateVerifiedAvailable,
		AccessStateLastVerified: "2026-05-07T10:00:00Z",
		KeyCustodyState:         KeyCustodyNotApplicable,
		RetentionLifecycle:      RetentionLifecycle{State: RetentionLifecycleActive, PolicyRef: "policy:retain-30d", ExpiresAt: "2026-06-06T10:00:00Z"},
		UnavailableReason:       UnavailableReasonNotAssessed,
	}
}
