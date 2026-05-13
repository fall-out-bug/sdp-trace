package prreview

func modelMismatchWithoutFallback(role ReviewRole, result ReviewerResult) bool {
	// modelMismatchWithoutFallback keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	requested := defaultString(role.RequestedModel, result.RequestedModel)
	observed := result.ObservedModel
	if modelIdentityMissing(requested) || modelIdentityMissing(observed) {
		return false
	}
	if requested == observed {
		return false
	}
	return fallbackMetadataMissing(result)
}
