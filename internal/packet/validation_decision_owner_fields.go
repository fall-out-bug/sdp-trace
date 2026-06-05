package packet

import "strings"

// Blank decisions are rejected before named-field validation; otherwise they
// could satisfy no required decision and still produce misleading field errors.
func (v *bundleValidator) validateDecisionOwner(owner DecisionOwner) string {
	decision := strings.TrimSpace(owner.Decision)
	if decision == "" {
		v.add("decision owner requires decision")
		return ""
	}
	v.validateNamedDecisionOwner(decision, owner)
	return decision
}

// Named owner validation mirrors row validation: owner identity, state, then
// reason for non-pass states.
func (v *bundleValidator) validateNamedDecisionOwner(decision string, owner DecisionOwner) {
	v.validateDecisionOwnerName(decision, owner)
	v.validateDecisionOwnerState(decision, owner)
	v.validateDecisionOwnerReason(decision, owner)
}

func (v *bundleValidator) validateDecisionOwnerName(decision string, owner DecisionOwner) {
	if strings.TrimSpace(owner.Owner) == "" {
		v.add("decision %s requires owner", decision)
	}
}

func (v *bundleValidator) validateDecisionOwnerState(decision string, owner DecisionOwner) {
	if !states[owner.State] {
		v.add("decision %s has unknown state %q", decision, owner.State)
	}
}

func (v *bundleValidator) validateDecisionOwnerReason(decision string, owner DecisionOwner) {
	if missingReasonStates[owner.State] && strings.TrimSpace(owner.Reason) == "" {
		v.add("decision %s state %s requires reason", decision, owner.State)
	}
}
