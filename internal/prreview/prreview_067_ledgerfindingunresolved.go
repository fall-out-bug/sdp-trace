package prreview

func ledgerFindingUnresolved(finding LedgerFinding) bool {
	return (finding.Severity == SeverityCritical || finding.Severity == SeverityMajor) && finding.Disposition == DispositionUnresolvedReviewBlocker
}
