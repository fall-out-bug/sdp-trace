package packet

func (v *bundleValidator) validateDecisionOwnerState(decision string, owner DecisionOwner) {
	// validateDecisionOwnerState keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if !states[owner.State] {

		v.add("decision %s has unknown state %q", decision, owner.State)
	}
}
