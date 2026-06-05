package prreview

// normalizeParsedReviewerOutput overlays runner-owned metadata onto parsed
// reviewer output before sanitizing it for storage.
func normalizeParsedReviewerOutput(parsed, base ReviewerResult, role ReviewRole) ReviewerResult {
	parsed.ReviewRunID = defaultString(parsed.ReviewRunID, base.ReviewRunID)
	parsed.Runner = defaultString(parsed.Runner, role.Runner)
	normalizeParsedReviewerModels(&parsed, role)
	attachParsedReviewerExecution(&parsed, base)
	normalizeParsedReviewerStatus(&parsed)
	return sanitizeReviewerResult(parsed)
}

func normalizeParsedReviewerModels(parsed *ReviewerResult, role ReviewRole) {
	parsed.RequestedModel = defaultString(parsed.RequestedModel, defaultString(role.RequestedModel, StateNotAssessed))
	parsed.ObservedModel = defaultString(parsed.ObservedModel, StateNotAssessed)
	parsed.ModelFamily = defaultString(parsed.ModelFamily, StateNotAssessed)
	parsed.ModelVersion = defaultString(parsed.ModelVersion, StateNotAssessed)
}

func attachParsedReviewerExecution(parsed *ReviewerResult, base ReviewerResult) {
	parsed.StartedAt = base.StartedAt
	parsed.EndedAt = base.EndedAt
	parsed.CommandDigest = base.CommandDigest
	parsed.PromptRef = base.PromptRef
}
