package forensic

func validProfileMappings() []ProfileMapping {
	// validProfileMappings keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.

	required := []string{
		RetentionModeSanitizedExcerpt,
		RetentionModeEncryptedRawRef,
		RetentionModeExternalArtifactRef,
	}

	authority := AuthorityRef{ActorID: "human:security-owner", VerificationState: AuthorityVerified}
	return []ProfileMapping{
		{EventFamily: "command_finished", RequiredRetentionModes: required, Critical: true, Authority: authority},
		{EventFamily: "test_output_observed", RequiredRetentionModes: required, Critical: true, Authority: authority},
	}
}

func validRedactionRules() []Rule {
	// validRedactionRules keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.

	return []Rule{
		{RuleID: "secret-token-v1", DetectorFamily: "secret", RuleVersion: "1.0.0", Action: RedactionActionApplyRule, RetentionMode: RetentionModeSanitizedExcerpt},
		{RuleID: "withhold-privacy-v1", DetectorFamily: "privacy", RuleVersion: "1.0.0", Action: RedactionActionWithhold, RetentionMode: RetentionModeNotAssessed},
	}
}
