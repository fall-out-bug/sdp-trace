package witness

func applyCIEnvelopeTrustDecision(record *Record, kind, runsRoot string, envelope EnvelopeInput) bool {
	// Validation first checks envelope self-consistency and artifact binding,
	// then separately binds the claimed run ID to discovered local runs.
	state := validateCIEnvelope(kind, envelope, record.RunArtifacts)
	if state.reason != "" {
		applyProfileState(record, state.status, state.scope, state.reason)
		return false
	}
	return setCIEnvelopeRunBindingState(record, runsRoot, envelope.CI.RunID)
}

func setCIEnvelopeRunBindingState(record *Record, runsRoot, witnessRunID string) bool {
	if runIDMatches(runsRoot, witnessRunID) {
		// A passing envelope profile requires both valid envelope states and a
		// run ID that resolves to a current discovered run artifact.
		record.Status = StatusPass
		record.TrustScope = TrustScopeCIWitnessed
		record.EstablishedTrustScope = TrustScopeCIWitnessed
		record.Reason = ReasonProfileVerified
		return true
	}
	record.ProfileStates.RunBindingState = stateFail
	applyProfileState(record, StatusFail, stateFail, ReasonRunMismatch)
	return false
}
