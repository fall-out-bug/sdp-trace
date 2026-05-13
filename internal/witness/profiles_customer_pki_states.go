package witness

func customerPKICannotVerify(record *Record, reason string) {
	// cannot_verify keeps the requested external profile visible while making
	// clear that no external trust scope was established.
	applyProfileState(record, StatusCannotVerify, stateCannotVerify, reason)
	record.ProfileStates = defaultProfileStates(stateCannotVerify, independenceExternal)
}

func customerPKIInputFail(record *Record, reason string) {
	// Unsafe PKI inputs fail the profile immediately because the evidence source
	// violates the witness safety contract.
	applyProfileState(record, StatusFail, stateFail, reason)
	record.ProfileStates = defaultProfileStates(stateFail, independenceExternal)
}

func customerPKIPassStates(policy CustomerPKIAuthorityPolicy) *ProfileStates {
	// These optimistic states are provisional: they become authoritative only
	// after validateCustomerPKIRecord completes every external-evidence check.
	return &ProfileStates{
		IdentityState:        statePass,
		SignerAuthorityState: statePass,
		FreshnessState:       statePass,
		ArtifactBindingState: statePass,
		SourceBindingState:   stateNotAssessed,
		RunBindingState:      statePass,
		PolicyBindingState:   statePass,
		IndependenceState:    independenceExternal,
		KeyCustodyState:      defaultString(policy.KeyCustodyState, "not_assessed"),
	}
}
