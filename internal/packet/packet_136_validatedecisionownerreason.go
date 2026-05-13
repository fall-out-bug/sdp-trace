package packet

import (
	"strings"
)

func (v *bundleValidator) validateDecisionOwnerReason(decision string, owner DecisionOwner) {
	// validateDecisionOwnerReason keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if missingReasonStates[owner.State] && strings.TrimSpace(owner.Reason) == "" {

		v.add("decision %s state %s requires reason", decision, owner.State)
	}
}
