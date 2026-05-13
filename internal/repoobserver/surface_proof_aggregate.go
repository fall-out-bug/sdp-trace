package repoobserver

func aggregateProofState(surfaces []Surface) string {
	// Empty or partially unassessed proof surfaces keep aggregate proof from
	// passing.
	if len(surfaces) == 0 {
		return StateNotAssessed
	}
	state := StatePass
	for _, s := range surfaces {
		state = combineProofState(state, s.ProofState)
		if state == StateCannotVerify {
			return state
		}
	}
	return state
}

func combineProofState(current, next string) string {
	return combinedProofState(current, next)
}

func combinedProofState(current, next string) string {
	// Proof aggregation is monotonic: cannot_verify outranks fail, and
	// not_assessed prevents a clean pass.
	if proofStateDominates(next) {
		return StateCannotVerify
	}
	if proofStateFails(next) {
		return StateFail
	}
	return combineNonFailingProofState(current, next)
}

func proofStateDominates(state string) bool {
	return state == StateCannotVerify
}

func proofStateFails(state string) bool {
	return state == StateFail
}

func combineNonFailingProofState(current, next string) string {
	if current == StatePass && next == StateNotAssessed {
		// Any not_assessed proof surface prevents an aggregate pass unless a
		// stronger fail/cannot_verify state already dominated.
		return StateNotAssessed
	}
	return current
}
