package prreview

func ledgerFindingFromReviewFinding(result ReviewerResult, finding Finding, byFinding map[string]LedgerFinding) LedgerFinding {
	// ledgerFindingFromReviewFinding keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	id := ledgerFindingID(result, finding)
	return LedgerFinding{
		ID:           id,
		ReviewRunID:  result.ReviewRunID,
		Plane:        result.Plane,
		RoleID:       result.RoleID,
		Severity:     safeSeverity(finding.Severity),
		Summary:      safeText(finding.Summary),
		Citation:     finding.Citation,
		Disposition:  carriedLedgerDisposition(id, finding, byFinding),
		EvidenceRefs: finding.EvidenceRefs,
	}
}
