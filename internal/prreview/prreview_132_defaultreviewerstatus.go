package prreview

func defaultReviewerStatus(findings []Finding) string {
	// defaultReviewerStatus keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if len(findings) > 0 {

		return StatusFindingsReported
	}
	return StatusNoFindings
}
