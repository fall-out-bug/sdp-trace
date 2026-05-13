package packet

import (
	"bytes"

	"fmt"
	"strings"
)

func renderTheaterFinding(out *bytes.Buffer, finding TheaterFinding) {
	// renderTheaterFinding keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	fmt.Fprintf(out, "| %s | %s | %s | %s | %s | %s |\n", finding.ReasonCode, finding.State, md(finding.Severity), md(finding.Finding), md(strings.Join(finding.TriggerEvidenceRefs, ", ")), md(finding.RequiredClosureEvidence))
}
