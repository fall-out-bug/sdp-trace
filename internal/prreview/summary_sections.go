package prreview

import (
	"fmt"
	"strings"
)

// writeSummaryPlanes renders plane status and safe next actions without
// expanding reviewer output into new proof claims.
func writeSummaryPlanes(b *strings.Builder, planes []PlaneResult) {
	if len(planes) == 0 {
		return
	}
	b.WriteString("\nPlanes\n")
	for _, plane := range planes {
		writeSummaryPlane(b, plane)
	}
}

// writeSummaryPlane keeps next actions sanitized because they may originate
// from reviewer or validation text.
func writeSummaryPlane(b *strings.Builder, plane PlaneResult) {
	fmt.Fprintf(b, "- %s: %s", plane.Plane, plane.Status)
	if plane.NextAction != "" {
		fmt.Fprintf(b, " next_action=%s", safeText(plane.NextAction))
	}
	b.WriteString("\n")
}

// writeSummaryFindings renders only safe finding summaries and their existing
// disposition state; it does not adjudicate or close findings.
func writeSummaryFindings(b *strings.Builder, findings []LedgerFinding) {
	if len(findings) == 0 {
		return
	}
	b.WriteString("\nFindings\n")
	for _, finding := range findings {
		fmt.Fprintf(b, "- %s [%s] %s (%s)\n", finding.ID, finding.Severity, safeText(finding.Summary), finding.Disposition)
	}
}
