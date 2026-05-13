package forensic

func profileSelectionCondition(input Input) Condition {
	// profileSelectionCondition keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	selection := input.Run.ProfileSelection
	if selection.SelectedProfile == "" {

		return Condition{ID: "profile_selection_accountable", State: StateNotAssessed, ReasonCode: "profile_selection_not_assessed", Reason: "profile selection accountability is not recorded", NextAction: "Record actor, profile, policy digest, and justification when policy requires it."}
	}
	if !profileSelectionVerified(selection, input.Policy.PolicyDigest) {
		return cannotVerify("profile_selection_accountable", "profile_selection_unverifiable", "forensic profile selection accountability cannot be verified", "Record accountable forensic profile selection evidence.")
	}
	return pass("profile_selection_accountable", "profile_selection_accountable", "forensic profile selection is accountable")
}

func profileSelectionVerified(selection ProfileSelection, policyDigest string) bool {
	// profileSelectionVerified keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.

	required := []bool{
		selection.SelectedProfile == ProfileForensicRetention,
		selection.RedactionPolicyDigest == policyDigest,
		selection.ActorID != "",
		selection.Justification != "",
		selection.AuthorityVerified,
	}
	for _, ok := range required {
		if !ok {

			return false
		}
	}
	return true
}
