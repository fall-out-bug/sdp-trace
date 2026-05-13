package prreview

import (
	"fmt"

	"strings"
)

func writeSummaryFindings(b *strings.Builder, findings []LedgerFinding) {
	// writeSummaryFindings keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if len(findings) > 0 {
		b.WriteString("\nFindings\n")
		for _, finding := range findings {

			fmt.Fprintf(b, "- %s [%s] %s (%s)\n", finding.ID, finding.Severity, safeText(finding.Summary), finding.Disposition)
		}
	}
}
