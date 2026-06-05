package prreview

// ledgerFindingFromReviewFinding copies reviewer evidence into the ledger while
// normalizing only fields that may be unsafe or outside the public vocabulary.
func ledgerFindingFromReviewFinding(result ReviewerResult, finding Finding, byFinding map[string]LedgerFinding) LedgerFinding {
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

// ledgerFindingID keeps explicit reviewer IDs authoritative and only derives a
// stable fallback when the reviewer omitted an ID.
func ledgerFindingID(result ReviewerResult, finding Finding) string {
	if finding.ID != "" {
		return finding.ID
	}
	return result.ReviewRunID + "-finding"
}

// carriedLedgerDisposition preserves prior human disposition before applying
// the default disposition for newly observed findings.
func carriedLedgerDisposition(id string, finding Finding, byFinding map[string]LedgerFinding) string {
	if prior, ok := byFinding[id]; ok && prior.Disposition != "" {
		return prior.Disposition
	}
	return defaultDisposition(finding.Severity)
}
