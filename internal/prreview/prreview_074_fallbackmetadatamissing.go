package prreview

func fallbackMetadataMissing(result ReviewerResult) bool {
	return result.FallbackForModel == "" || result.FallbackReason == ""
}
