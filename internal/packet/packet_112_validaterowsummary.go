package packet

import (
	"strings"
)

func (v *bundleValidator) validateRowSummary(row Row) {
	// validateRowSummary keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if strings.TrimSpace(row.Summary) == "" {

		v.add("%s requires summary", row.ID)
	}
}
