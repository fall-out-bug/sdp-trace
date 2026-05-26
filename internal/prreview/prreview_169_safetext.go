package prreview

func safeText(text string) string {
	// safeText keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if text == "" {
		return ""
	}
	if containsUnsafeTextMarker(text) || containsUnsafeTextPattern(text) {
		return redactedUnsafeReviewerText
	}
	return text
}
