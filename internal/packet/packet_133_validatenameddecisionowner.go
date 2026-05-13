package packet

func (v *bundleValidator) validateNamedDecisionOwner(decision string, owner DecisionOwner) {
	// validateNamedDecisionOwner keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	v.validateDecisionOwnerName(decision, owner)
	v.validateDecisionOwnerState(decision, owner)
	v.validateDecisionOwnerReason(decision, owner)
}
