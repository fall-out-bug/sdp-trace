package prreview

func normalizeParsedReviewerStatus(parsed *ReviewerResult) {
	if parsed.Status == "" {
		parsed.Status = defaultReviewerStatus(parsed.Findings)
	}
}

func defaultReviewerStatus(findings []Finding) string {
	if len(findings) > 0 {
		return StatusFindingsReported
	}
	return StatusNoFindings
}
