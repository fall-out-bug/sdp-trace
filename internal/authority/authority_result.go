package authority

import (
	"strings"
)

func authorityResult(pkg Package, actions []ObservedAction, bindings []EvidenceBinding, evaluations []AuthorityEvaluation, envState string) Result {
	// authorityResult keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.

	result := Result{
		SchemaVersion:            ResultSchemaVersion,
		SelectedProfile:          Profile,
		SelectedPolicyID:         strings.TrimSpace(pkg.SelectedPolicyID),
		AuthorityEvaluationState: aggregateState(evaluations, envState),
		Evaluations:              evaluations,
		BindingEvaluations:       bindings,
		SourceCoverage:           sourceCoverage(actions),
	}

	result.Reasons = resultReasons(result)
	result.NextActions = nextActions(result)
	return result
}
