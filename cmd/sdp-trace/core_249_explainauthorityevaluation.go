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
