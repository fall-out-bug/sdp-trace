package prreview

func sanitizeReviewerResult(result ReviewerResult) ReviewerResult {
	// Reviewer output is untrusted even when structurally valid JSON.
	// Every user-visible text channel is redacted before runs, ledgers, and
	// summaries can persist or upload it.
	result.Reason = safeText(result.Reason)
	for i := range result.Findings {
		result.Findings[i].Summary = safeText(result.Findings[i].Summary)
		result.Findings[i].Rationale = safeText(result.Findings[i].Rationale)
		result.Findings[i].SuggestedFix = safeText(result.Findings[i].SuggestedFix)
		result.Findings[i].Question = safeText(result.Findings[i].Question)
		for j := range result.Findings[i].EvidenceRefs {
			result.Findings[i].EvidenceRefs[j] = safeEvidenceRef(result.Findings[i].EvidenceRefs[j])
		}
	}
	return result
}

func safeEvidenceRef(ref string) string {
	values := map[bool]string{true: redactedUnsafeReviewerText, false: safeID(ref)}
	return values[unsafeEvidenceRef(ref)]
}

func unsafeEvidenceRef(ref string) bool {
	return containsUnsafeTextMarker(ref) || containsUnsafeTextPattern(ref)
}
