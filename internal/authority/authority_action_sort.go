package authority

import (
	"fmt"
	"sort"
)

func sortedObservedActions(actions []ObservedAction) []ObservedAction {
	// sortedObservedActions keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.

	out := append([]ObservedAction(nil), actions...)
	sort.Slice(out, func(i, j int) bool { return out[i].EventID < out[j].EventID })
	return out
}

func evaluateActions(pkg Package, env AuthorityEnvelope, envState, envReason string, actions []ObservedAction, bindings []EvidenceBinding) []AuthorityEvaluation {
	// evaluateActions keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.

	bindingByEvent := bindingStatesByEvent(bindings)
	resolution := evidenceResolutionIndex(pkg.EvidenceResolution)
	evaluations := make([]AuthorityEvaluation, 0, len(actions))
	for i, action := range actions {

		evaluationID := fmt.Sprintf("authority-evaluation-%03d", i+1)
		evaluations = append(evaluations, evaluateAction(evaluationID, pkg.SelectedPolicyID, env, envState, envReason, action, bindingByEvent[action.EventID], resolution))
	}
	return evaluations
}
