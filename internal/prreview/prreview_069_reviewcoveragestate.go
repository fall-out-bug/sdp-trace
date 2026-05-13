package prreview

func reviewCoverageState(required map[string]bool, usableCount int, cannotVerify, unresolved bool) string {
	// reviewCoverageState keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if cannotVerify {
		return CoverageCannotVerify
	}
	if noReviewCoverage(required, usableCount) {
		return CoverageNotAssessed
	}
	return assessedReviewCoverageState(required, usableCount, unresolved)
}
