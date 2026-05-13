package ciartifact

func reasons(families []FamilyObservation, index ArtifactIndexResult, safety OutputSafetyResult, identityCannotVerify bool) []string {
	// Reasons are derived from the recorded family states so prose follows machine
	// evidence instead of becoming independent authority.
	set := map[string]bool{}

	addFamilyReasons(set, families)
	addVisibleReason(set, index.Result, index.ReasonCode, index.Reason)
	addVisibleReason(set, safety.State, safety.ReasonCode, safety.Reason)
	addIdentityReason(set, identityCannotVerify)
	return sortedKeys(set)
}

func addFamilyReasons(set map[string]bool, families []FamilyObservation) {
	// addFamilyReasons keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	for _, family := range families {

		addVisibleReason(set, family.FamilyState, family.ReasonCode, family.Reason)
	}
}

func addVisibleReason(set map[string]bool, state, code, reason string) {
	// addVisibleReason keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if familyReasonVisible(state) {

		set[code+": "+reason] = true
	}
}

func addIdentityReason(set map[string]bool, identityCannotVerify bool) {
	if identityCannotVerify {
		set["unsafe_identity_metadata: selected source or run identity contained unsafe or unsupported metadata"] = true
	}
}

func familyReasonVisible(state string) bool {
	return state != StatePass && state != StateNotAssessed
}
