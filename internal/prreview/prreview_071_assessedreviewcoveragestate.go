package prreview

func assessedReviewCoverageState(required map[string]bool, usableCount int, unresolved bool) string {
	// assessedReviewCoverageState keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if usableCount < len(required) {

		return CoveragePartial
	}
	if unresolved {

		return CoverageUnresolved
	}
	return CoverageSatisfied
}
