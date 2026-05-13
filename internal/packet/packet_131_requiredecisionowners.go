package packet

func (v *bundleValidator) requireDecisionOwners(owners map[string]DecisionOwner) {
	// requireDecisionOwners keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	for _, decision := range requiredDecisions {
		if owners[decision].Decision == "" {

			v.add("missing decision owner %q", decision)
		}
	}
}
