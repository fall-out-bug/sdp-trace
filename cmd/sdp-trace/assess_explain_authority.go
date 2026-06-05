package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/authority"
)

func explainAuthorityEvaluation(result authority.Result, stdout io.Writer) int {
	// Authority explanations separate action policy results from evidence
	// binding results so provenance failures stay visible.
	fmt.Fprintf(stdout, "Selected profile: %s\n", result.SelectedProfile)
	fmt.Fprintf(stdout, "Authority evaluation: %s\n", result.AuthorityEvaluationState)
	fmt.Fprintf(stdout, "Selected policy: %s\n", result.SelectedPolicyID)
	explainAuthorityActionEvaluations(result.Evaluations, stdout)
	explainAuthorityBindingEvaluations(result.BindingEvaluations, stdout)
	explainReasons(result.Reasons, stdout)
	explainNextActions(result.NextActions, stdout)
	return 0
}

func explainAuthorityActionEvaluations(evaluations []authority.AuthorityEvaluation, stdout io.Writer) {
	for _, eval := range evaluations {
		// Attribution fields are printed independently; a matched action rule
		// does not imply actor, tool, or model provenance is complete.
		fmt.Fprintf(stdout, "Observed action %s: %s (%s)\n", eval.EventID, eval.State, eval.ReasonCode)
		fmt.Fprintf(stdout, "  Actor attribution: %s\n", eval.ActorAttribution)
		fmt.Fprintf(stdout, "  Tool attribution: %s\n", eval.ToolAttribution)
		fmt.Fprintf(stdout, "  Model attribution: %s\n", eval.ModelAttribution)
		if eval.MatchedRuleRef != "" {
			fmt.Fprintf(stdout, "  Matched rule: %s\n", eval.MatchedRuleRef)
		}
	}
}

func explainAuthorityBindingEvaluations(bindings []authority.EvidenceBinding, stdout io.Writer) {
	for _, binding := range bindings {
		// Binding lines keep evidence-binding state separate from action state
		// so provenance gaps remain visible in CLI output.
		fmt.Fprintf(stdout, "Binding %s: %s (%s)\n", binding.BindingID, binding.BindingState, binding.ReasonCode)
	}
}
