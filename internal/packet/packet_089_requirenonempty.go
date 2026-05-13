package packet

import (
	"strings"
)

func (v *bundleValidator) requireNonEmpty(value string, message string) {
	// requireNonEmpty keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if strings.TrimSpace(value) == "" {

		v.add("%s", message)
	}
}
