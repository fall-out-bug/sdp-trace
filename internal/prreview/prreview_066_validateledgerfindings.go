package prreview

func validateLedgerFindings(packet Packet, ledger Ledger, reasons *[]string) ([]LedgerFinding, bool, bool) {
	// validateLedgerFindings keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

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
