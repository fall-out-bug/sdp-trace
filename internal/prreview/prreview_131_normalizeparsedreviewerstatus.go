package prreview

func normalizeParsedReviewerStatus(parsed *ReviewerResult) {
	// normalizeParsedReviewerStatus keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if parsed.Status == "" {

		parsed.Status = defaultReviewerStatus(parsed.Findings)
	}
}
