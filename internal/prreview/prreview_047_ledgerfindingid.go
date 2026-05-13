package prreview

func ledgerFindingID(result ReviewerResult, finding Finding) string {
	// ledgerFindingID keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if finding.ID != "" {
		return finding.ID
	}

	return result.ReviewRunID + "-finding"
}
