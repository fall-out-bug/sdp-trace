package packet

import (
	"strings"
)

func (v *bundleValidator) validateRowReason(row Row) {
	// validateRowReason keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	if missingReasonStates[row.State] && strings.TrimSpace(row.Reason) == "" {
		v.add("%s state %s requires reason", row.ID, row.State)
	}
}
