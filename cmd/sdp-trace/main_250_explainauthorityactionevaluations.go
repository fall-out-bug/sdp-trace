package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/authority"
)

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
