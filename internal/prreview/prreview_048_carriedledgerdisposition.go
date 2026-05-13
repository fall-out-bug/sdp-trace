package prreview

func carriedLedgerDisposition(id string, finding Finding, byFinding map[string]LedgerFinding) string {
	// carriedLedgerDisposition keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if prior, ok := byFinding[id]; ok && prior.Disposition != "" {

		return prior.Disposition
	}
	return defaultDisposition(finding.Severity)
}
