package forensic

func validRunEvidence(policyDigest string) RunEvidence {
	// validRunEvidence keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.

	return RunEvidence{
		RunID:                 "forensic-run-1",
		SelectedProfile:       ProfileForensicRetention,
		RedactionPolicyDigest: policyDigest,
		ProfileSelection:      validProfileSelection(policyDigest),
		Events:                validEventRetentions(policyDigest),
	}
}
func validProfileSelection(policyDigest string) ProfileSelection {
	// validProfileSelection keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.

	return ProfileSelection{
		ActorID:                 "human:security-owner",
		SelectedProfile:         ProfileForensicRetention,
		RedactionPolicyDigest:   policyDigest,
		Justification:           "incident review",
		AuthorityVerified:       true,
		SelectionEvidenceDigest: "2222222222222222222222222222222222222222222222222222222222222222",
	}
}
