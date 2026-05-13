package authority

import (
	"sort"
)

func evaluateBindings(inputs []EvidenceBindingInput, actions []ObservedAction) []EvidenceBinding {
	// evaluateBindings keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	actionIDs := map[string]bool{}
	for _, action := range actions {

		actionIDs[action.EventID] = true
	}
	out := make([]EvidenceBinding, 0, len(inputs))
	for _, input := range inputs {
		out = append(out, evaluateBinding(input, actionIDs))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].BindingID < out[j].BindingID })
	return out
}

func evaluateBinding(input EvidenceBindingInput, actionIDs map[string]bool) EvidenceBinding {
	// evaluateBinding keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.

	state, reason := bindingStateAndReason(input, actionIDs)
	return EvidenceBinding{
		BindingID:     input.BindingID,
		LeftEventID:   input.LeftEventID,
		RightEventID:  input.RightEventID,
		BindingType:   input.BindingType,
		BindingState:  state,
		MatchedFields: input.MatchedFields,
		EvidenceRef:   input.EvidenceRef,
		ReasonCode:    reason,
	}
}
