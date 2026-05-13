package prreview

func normalizeParsedReviewerModels(parsed *ReviewerResult, role ReviewRole) {
	// normalizeParsedReviewerModels keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	parsed.RequestedModel = defaultString(parsed.RequestedModel, defaultString(role.RequestedModel, StateNotAssessed))
	parsed.ObservedModel = defaultString(parsed.ObservedModel, StateNotAssessed)
	parsed.ModelFamily = defaultString(parsed.ModelFamily, StateNotAssessed)
	parsed.ModelVersion = defaultString(parsed.ModelVersion, StateNotAssessed)
}
