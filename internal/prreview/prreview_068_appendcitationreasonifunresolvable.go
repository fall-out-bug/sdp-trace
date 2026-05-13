package prreview

func appendCitationReasonIfUnresolvable(packet Packet, finding LedgerFinding, reasons *[]string) bool {
	// appendCitationReasonIfUnresolvable keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if citationResolvable(packet, finding.Citation) {
		return false
	}

	*reasons = append(*reasons, "finding_citation_cannot_verify")
	return true
}
