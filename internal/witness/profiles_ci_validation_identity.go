package witness

func validateCIEnvelopeIdentityAuthorityStates(states ProfileStates) profileDecision {
	// Missing identity or signer authority means the envelope never reached a
	// replayable evidence boundary, so the verdict stays cannot_verify.
	if states.IdentityState != statePass {
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonMissingIdentity}
	}
	if states.SignerAuthorityState != statePass {
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonMissingSigner}
	}
	return profileDecision{}
}

func validateCIEnvelopeFreshnessState(states ProfileStates) profileDecision {
	// Explicit stale freshness is contradictory evidence and fails; any other
	// non-pass freshness state is missing evidence.
	if states.FreshnessState == stateFail {
		return profileDecision{StatusFail, stateFail, ReasonStaleFreshness}
	}
	if states.FreshnessState != statePass {
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonMissingFreshness}
	}
	return profileDecision{}
}
