package packet

import (
	"strings"
)

func (v *bundleValidator) validateDecisionOwnerName(decision string, owner DecisionOwner) {
	// validateDecisionOwnerName keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if strings.TrimSpace(owner.Owner) == "" {

		v.add("decision %s requires owner", decision)
	}
}
