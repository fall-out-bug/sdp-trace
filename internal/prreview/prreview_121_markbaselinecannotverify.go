package prreview

func markBaselineCannotVerify(result *ReviewerResult) {
	result.Status = StatusCannotVerify
	result.Reason = "working_tree_baseline_cannot_verify"
}
