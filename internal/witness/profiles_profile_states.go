package witness

func defaultProfileStates(state, independence string) *ProfileStates {
	// Defaults make unassessed profile state explicit; callers must override
	// individual states when they have evidence for a narrower verdict.
	return &ProfileStates{
		IdentityState:        state,
		SignerAuthorityState: state,
		FreshnessState:       state,
		ArtifactBindingState: state,
		SourceBindingState:   state,
		RunBindingState:      state,
		PolicyBindingState:   state,
		IndependenceState:    independence,
	}
}
