package prreview

func safeSeverity(severity string) string {
	// safeSeverity keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	switch severity {
	case SeverityCritical, SeverityMajor, SeverityMinor, SeverityInformational:

		return severity
	default:
		return SeverityInformational
	}
}
