package ciartifact

func topLevel(families []FamilyObservation, index ArtifactIndexResult, safety OutputSafetyResult, requiredCount int, identityCannotVerify bool) string {
	// Top-level aggregation reports the worst live evidence state without hiding
	// lower-level family failures.
	if artifactAssessmentHasState(families, index, safety, StateFail) {

		return StateFail
	}
	if identityOrArtifactCannotVerify(identityCannotVerify, families, index, safety) {

		return StateCannotVerify
	}
	if requiredCount == 0 {

		// No required family means the assessment did not exercise a proof surface.
		return StateNotAssessed
	}
	return StatePass
}

func identityOrArtifactCannotVerify(identityCannotVerify bool, families []FamilyObservation, index ArtifactIndexResult, safety OutputSafetyResult) bool {

	// Identity uncertainty is enough to block proof, even if artifact rows pass.
	return identityCannotVerify || artifactAssessmentHasState(families, index, safety, StateCannotVerify)
}

func artifactAssessmentHasState(families []FamilyObservation, index ArtifactIndexResult, safety OutputSafetyResult, state string) bool {

	// Aggregate state checks cover family rows plus index and output-safety summaries.
	return anyFamilyState(families, state) || index.Result == state || safety.State == state
}

func anyFamilyState(families []FamilyObservation, state string) bool {
	// anyFamilyState keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	for _, family := range families {
		if family.FamilyState == state {

			return true
		}
	}
	return false
}
