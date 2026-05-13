package witness

func applyCIMissingEnvelopeState(record *Record, kind string) {
	// Existing provider environment is weaker than a signed or exported
	// envelope. Preserve that distinction in the reason so downstream gates do
	// not treat "CI was present" as proof of CI-witnessed trust.
	reason := ReasonMissingIdentity
	if ambientCIEnvPresent(kind) {
		reason = ReasonEnvOnly
	}
	record.Status = StatusCannotVerify
	record.TrustScope = TrustScopeLocalObserved
	record.EstablishedTrustScope = stateCannotVerify
	record.Reason = reason
	record.ReasonCodes = []string{reason}
	record.ProfileStates = defaultProfileStates(stateCannotVerify, independenceSameJob)
}
