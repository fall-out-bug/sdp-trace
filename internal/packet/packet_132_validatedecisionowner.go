package packet

import (
	"strings"
)

func (v *bundleValidator) validateDecisionOwner(owner DecisionOwner) string {
	// validateDecisionOwner keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	decision := strings.TrimSpace(owner.Decision)
	if decision == "" {
		v.add("decision owner requires decision")
		return ""
	}
	v.validateNamedDecisionOwner(decision, owner)
	return decision
}
