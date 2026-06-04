package prreview

// validateLedgerFindings returns the safe ledger findings that validation can
// render, plus coverage signals derived from unresolved blockers and citations.
func validateLedgerFindings(packet Packet, ledger Ledger, reasons *[]string) ([]LedgerFinding, bool, bool) {
	unresolved := false
	cannotVerify := false
	safeFindings := make([]LedgerFinding, 0, len(ledger.Findings))
	for _, finding := range ledger.Findings {
		finding.Summary = safeText(finding.Summary)
		unresolved = unresolved || ledgerFindingUnresolved(finding)
		cannotVerify = cannotVerify || appendCitationReasonIfUnresolvable(packet, finding, reasons)
		safeFindings = append(safeFindings, finding)
	}
	return safeFindings, unresolved, cannotVerify
}

// ledgerFindingUnresolved treats only critical and major unresolved blockers as
// coverage-blocking; lower severities still remain visible in validation output.
func ledgerFindingUnresolved(finding LedgerFinding) bool {
	return (finding.Severity == SeverityCritical || finding.Severity == SeverityMajor) && finding.Disposition == DispositionUnresolvedReviewBlocker
}

// appendCitationReasonIfUnresolvable records the shared validation reason for
// ledger findings whose citation cannot be resolved against the packet refs.
func appendCitationReasonIfUnresolvable(packet Packet, finding LedgerFinding, reasons *[]string) bool {
	if citationResolvable(packet, finding.Citation) {
		return false
	}

	*reasons = append(*reasons, "finding_citation_cannot_verify")
	return true
}
