package witness

func validateCIEnvelopeStates(states ProfileStates) profileDecision {
	// State validation preserves the verifier's evidence taxonomy: missing
	// evidence remains cannot_verify while contradictory bindings fail.
	// The first matching state wins so the output reason stays deterministic and
	// mirrors the profile repair order.
	// Independence is evaluated last because it cannot rescue missing identity,
	// signer, freshness, binding, or artifact evidence.
	if decision := validateCIEnvelopeEvidenceStates(states); decision.reason != "" {
		return decision
	}
	if decision := validateCIEnvelopeBindingStates(states); decision.reason != "" {
		return decision
	}
	if !ciEnvelopeIndependent(states) {
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonEnvOnly}
	}
	return profileDecision{}
}

func validateCIEnvelopeEvidenceStates(states ProfileStates) profileDecision {
	// Evidence states are checked before binding states because missing identity,
	// signer authority, or freshness means the envelope cannot yet be trusted
	// enough to interpret source, run, policy, or artifact claims.
	if decision := validateCIEnvelopeIdentityAuthorityStates(states); decision.reason != "" {
		return decision
	}
	return validateCIEnvelopeFreshnessState(states)
}
