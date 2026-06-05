package prreview

// modelMismatchWithoutFallback flags an observed reviewer model that differs
// from the requested model unless fallback provenance explains the change.
func modelMismatchWithoutFallback(role ReviewRole, result ReviewerResult) bool {
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

// modelIdentityMissing keeps unknown and explicitly not_assessed model values
// from becoming false-positive identity mismatches.
func modelIdentityMissing(model string) bool {
	return model == "" || model == StateNotAssessed
}

// fallbackMetadataMissing requires both the original requested model and a
// reason before a reviewer result can rely on fallback provenance.
func fallbackMetadataMissing(result ReviewerResult) bool {
	return result.FallbackForModel == "" || result.FallbackReason == ""
}
