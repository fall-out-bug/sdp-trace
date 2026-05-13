package prreview

func normalizeParsedReviewerOutput(parsed, base ReviewerResult, role ReviewRole) ReviewerResult {
	// normalizeParsedReviewerOutput keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	parsed.ReviewRunID = defaultString(parsed.ReviewRunID, base.ReviewRunID)
	parsed.Runner = defaultString(parsed.Runner, role.Runner)
	normalizeParsedReviewerModels(&parsed, role)
	attachParsedReviewerExecution(&parsed, base)
	normalizeParsedReviewerStatus(&parsed)
	return parsed
}
