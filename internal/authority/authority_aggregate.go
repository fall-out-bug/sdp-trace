package authority

func aggregateState(evaluations []AuthorityEvaluation, envState string) string {
	// aggregateState keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	if envState == StateCannotVerify {

		return StateCannotVerify
	}
	if len(evaluations) == 0 {

		return StateNotAssessed
	}
	rank := highestEvaluationStateRank(evaluations)
	if rank < 0 {
		return StateCannotVerify
	}
	return aggregateStateByRank[rank]
}

func highestEvaluationStateRank(evaluations []AuthorityEvaluation) int {
	// highestEvaluationStateRank keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	highestRank := -1
	for _, eval := range evaluations {
		rank, ok := aggregateStatePriority[eval.State]
		if ok && rank > highestRank {

			highestRank = rank
		}
	}
	return highestRank
}
