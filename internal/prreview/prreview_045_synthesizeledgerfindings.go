package prreview

func synthesizeLedgerFindings(runs RunSet, byFinding map[string]LedgerFinding) []LedgerFinding {
	// synthesizeLedgerFindings keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	findings := []LedgerFinding{}
	for _, result := range runs.Results {
		for _, finding := range result.Findings {
			findings = append(findings, ledgerFindingFromReviewFinding(result, finding, byFinding))
		}
	}
	return findings
}
