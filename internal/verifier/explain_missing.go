package verifier

import (
	"fmt"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func appendMissingEvidence(lines []string, rows []trace.MissingEvidenceRow) []string {
	if len(rows) == 0 {
		return lines
	}
	// Missing evidence is rendered only when the verifier produced table rows;
	// absence of this section is not an independent green claim.
	lines = append(lines, "missing_evidence:")
	for _, row := range rows {
		lines = append(lines, missingEvidenceLine(row))
	}
	return lines
}

func missingEvidenceLine(row trace.MissingEvidenceRow) string {
	// Keep each row compact so explain output remains copyable into reviews
	// without losing the expected event, observed state, or replay reason.
	return fmt.Sprintf(" - %s: %s (%s)", row.ExpectedEvent, row.ObservedState, row.Reason)
}
