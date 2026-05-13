package witness

func validateCIEnvelopeBindingStates(states ProfileStates) profileDecision {
	// Binding checks run only after core evidence is present; at that point
	// source/run mismatches are failures, while policy/artifact gaps remain
	// cannot_verify because the required comparison evidence is incomplete.
	if decision := validateCIEnvelopeRunSourceStates(states); decision.reason != "" {
		return decision
	}
	return validateCIEnvelopePolicyArtifactStates(states)
}

func validateCIEnvelopeRunSourceStates(states ProfileStates) profileDecision {
	// Source and run mismatches contradict the claimed execution context, so
	// they fail rather than remaining open as missing evidence.
	if states.SourceBindingState != statePass {
		return profileDecision{StatusFail, stateFail, ReasonSourceMismatch}
	}
	if states.RunBindingState != statePass {
		return profileDecision{StatusFail, stateFail, ReasonRunMismatch}
	}
	return profileDecision{}
}

func validateCIEnvelopePolicyArtifactStates(states ProfileStates) profileDecision {
	// Policy and artifact binding gaps leave the envelope unverifiable because
	// the local verifier cannot prove the claimed evidence set.
	if states.PolicyBindingState != statePass {
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonPolicyMissing}
	}
	if states.ArtifactBindingState != statePass {
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonMissingArtifact}
	}
	return profileDecision{}
}

func ciEnvelopeIndependent(states ProfileStates) bool {
	return states.IndependenceState == independenceCIJob || states.IndependenceState == independenceExternal
}
