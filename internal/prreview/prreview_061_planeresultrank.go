package prreview

func planeResultRank(result PlaneResult) int {
	// Prefer retained, usable reviewer output over failed attempts. Findings
	// outrank no-findings so a later clean run cannot hide reported blockers.
	if result.Usable {
		return usablePlaneResultRank(result.Status)
	}
	return nonUsablePlaneResultRank(result.Status)
}
