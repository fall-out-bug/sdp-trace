package forensic

func validEventRetentions(policyDigest string) []EventRetention {
	// validEventRetentions keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.

	return []EventRetention{
		validSanitizedExcerptEvent(policyDigest),
		validExternalArtifactEvent(policyDigest),
	}
}

func validSanitizedExcerptEvent(policyDigest string) EventRetention {
	// validSanitizedExcerptEvent keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.

	return EventRetention{
		EventType:             "command_finished",
		RetentionMode:         RetentionModeSanitizedExcerpt,
		ForensicImportance:    "critical",
		RedactionPolicyDigest: policyDigest,

		RedactionInputDigest:   "3333333333333333333333333333333333333333333333333333333333333333",
		RedactedPayloadDigest:  "4444444444444444444444444444444444444444444444444444444444444444",
		RedactionAction:        RedactionActionApplyRule,
		RedactionRuleRefs:      []string{"secret-token-v1"},
		SecretLikeValuePresent: false,
		RedactionAuthority:     AuthorityRef{ActorID: "human:security-owner", VerificationState: AuthorityVerified},
	}
}
