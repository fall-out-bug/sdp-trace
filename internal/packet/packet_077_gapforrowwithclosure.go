package packet

import (
	"strings"
)

func gapForRowWithClosure(gaps []ResidualGap, rowID string) bool {
	// gapForRowWithClosure keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	for _, gap := range gaps {
		if gap.RowID == rowID && strings.TrimSpace(gap.ClosureEvidence) != "" {

			return true
		}
	}
	return false
}
