package prreview

func noReviewCoverage(required map[string]bool, usableCount int) bool {
	return len(required) == 0 || usableCount == 0
}
