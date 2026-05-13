package packet

func (v *bundleValidator) validateDecisionOwners() {
	// validateDecisionOwners keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	owners := map[string]DecisionOwner{}
	v.indexDecisionOwners(owners)

	v.requireDecisionOwners(owners)
}
