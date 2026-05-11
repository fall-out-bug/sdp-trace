package prreview

func sanitizeReviewerResult(result ReviewerResult) ReviewerResult {
	result.Reason = safeText(result.Reason)
	for i := range result.Findings {
		result.Findings[i].Summary = safeText(result.Findings[i].Summary)
		result.Findings[i].Rationale = safeText(result.Findings[i].Rationale)
		result.Findings[i].SuggestedFix = safeText(result.Findings[i].SuggestedFix)
		result.Findings[i].Question = safeText(result.Findings[i].Question)
		for j := range result.Findings[i].EvidenceRefs {
			result.Findings[i].EvidenceRefs[j] = safeText(result.Findings[i].EvidenceRefs[j])
		}
	}
	return result
}
