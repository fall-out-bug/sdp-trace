package prreview

func existingFindings(existing *Ledger) map[string]LedgerFinding {
	// existingFindings keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	byFinding := map[string]LedgerFinding{}
	if existing != nil {
		for _, finding := range existing.Findings {
			byFinding[finding.ID] = finding
		}
	}
	return byFinding
}
