package prreview

func usablePlaneResultRank(status string) int {
	// A usable findings result must dominate a usable no-findings result so a
	// later clean retry cannot erase blockers that remain in the ledger.
	if status == StatusFindingsReported {
		return 4
	}
	return 3
}
