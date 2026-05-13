package prreview

func defaultDisposition(severity string) string {
	// defaultDisposition keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	switch safeSeverity(severity) {
	case SeverityCritical, SeverityMajor:
		return DispositionUnresolvedReviewBlocker
	default:
		return DispositionDeferredNotAssessed
	}
}
