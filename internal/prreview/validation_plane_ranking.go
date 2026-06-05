package prreview

func planeResultRank(result PlaneResult) int {
	// Prefer retained, usable reviewer output over failed attempts. Findings
	// outrank no-findings so a later clean run cannot hide reported blockers.
	if result.Usable {
		return usablePlaneResultRank(result.Status)
	}
	return nonUsablePlaneResultRank(result.Status)
}

func usablePlaneResultRank(status string) int {
	// A usable findings result must dominate a usable no-findings result so a
	// later clean retry cannot erase blockers that remain in the ledger.
	if status == StatusFindingsReported {
		return 4
	}
	return 3
}

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
