package prreview

func attachParsedReviewerExecution(parsed *ReviewerResult, base ReviewerResult) {
	// attachParsedReviewerExecution keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	parsed.StartedAt = base.StartedAt
	parsed.EndedAt = base.EndedAt
	parsed.CommandDigest = base.CommandDigest
	parsed.PromptRef = base.PromptRef
}
