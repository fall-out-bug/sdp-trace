package packet

import (
	"strings"
)

func (v *bundleValidator) validateRowOwner(row Row) {
	// validateRowOwner keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if strings.TrimSpace(row.Owner) == "" {

		v.add("%s requires owner", row.ID)
	}
}
