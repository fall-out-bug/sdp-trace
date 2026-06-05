package prreview

func reviewCoverageState(required map[string]bool, usableCount int, cannotVerify, unresolved bool) string {
	// Cannot-verify evidence dominates coverage completeness: stale, malformed,
	// or unreplayable review evidence cannot be converted into an approval.
	if cannotVerify {
		return CoverageCannotVerify
	}
	if noReviewCoverage(required, usableCount) {
		return CoverageNotAssessed
	}
	return assessedReviewCoverageState(required, usableCount, unresolved)
}

func noReviewCoverage(required map[string]bool, usableCount int) bool {
	return len(required) == 0 || usableCount == 0
}

func assessedReviewCoverageState(required map[string]bool, usableCount int, unresolved bool) string {
	if usableCount < len(required) {
		return CoveragePartial
	}
	if unresolved {
		return CoverageUnresolved
	}
	return CoverageSatisfied
}
