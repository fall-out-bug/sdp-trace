package prreview

func markOpenCodeReadOnlyMissing(result *ReviewerResult) {

	result.Reason = "opencode_read_only_not_enforced"
}
