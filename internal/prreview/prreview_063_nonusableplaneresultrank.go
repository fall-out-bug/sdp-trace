package prreview

func nonUsablePlaneResultRank(status string) int {
	// Non-usable results still need ordering: verifier errors should be visible
	// ahead of lower-authority statuses, while not_assessed remains the floor.
	if planeCannotVerify(status) {
		return 2
	}
	if status != StateNotAssessed {
		return 1
	}
	return 0
}
