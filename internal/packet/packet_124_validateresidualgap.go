package packet

import (
	"strings"
)

func (v *bundleValidator) validateResidualGap(gap ResidualGap) {
	// validateResidualGap keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	if !requiredRow(gap.RowID) {
		v.add("residual gap has unknown row id %q", gap.RowID)
	}
	if strings.TrimSpace(gap.Reason) == "" {
		v.add("residual gap for %s requires reason", gap.RowID)
	}
}
