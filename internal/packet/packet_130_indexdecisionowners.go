package packet

func (v *bundleValidator) indexDecisionOwners(owners map[string]DecisionOwner) {
	// indexDecisionOwners keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	for _, owner := range v.bundle.Packet.DecisionOwners {
		if decision := v.validateDecisionOwner(owner); decision != "" {

			owners[decision] = owner
		}
	}
}
