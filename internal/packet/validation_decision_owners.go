package packet

// Decision-owner validation records who must make downstream human decisions;
// it deliberately does not turn this packet into an approval.
func (v *bundleValidator) validateDecisionOwners() {
	owners := map[string]DecisionOwner{}
	v.indexDecisionOwners(owners)
	v.requireDecisionOwners(owners)
}

// Duplicate decisions preserve pre-slice behavior: after validating each row,
// the last valid owner wins for required-decision presence checks.
func (v *bundleValidator) indexDecisionOwners(owners map[string]DecisionOwner) {
	for _, owner := range v.bundle.Packet.DecisionOwners {
		if decision := v.validateDecisionOwner(owner); decision != "" {
			// Preserve input-order overwrite semantics: the last valid owner for a
			// decision is the owner used by required-decision presence checks.
			owners[decision] = owner
		}
	}
}

// Required decisions are checked in contract order so diagnostics stay stable
// and reviewers can compare missing ownership without sorting noise.
func (v *bundleValidator) requireDecisionOwners(owners map[string]DecisionOwner) {
	for _, decision := range requiredDecisions {
		if owners[decision].Decision == "" {
			v.add("missing decision owner %q", decision)
		}
	}
}
