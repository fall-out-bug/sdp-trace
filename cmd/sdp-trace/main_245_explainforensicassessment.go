package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/forensic"
)

func explainForensicAssessment(result forensic.AssessmentResult, stdout io.Writer) int {
	// Forensic explanations restate retention findings from the artifact so
	// missing raw references stay represented as recorded assessment state.
	fmt.Fprintf(stdout, "Selected profile: %s\n", result.SelectedProfile)
	fmt.Fprintf(stdout, "Forensic retention assessment: %s\n", result.ForensicRetentionAssessment)
	fmt.Fprintf(stdout, "Trust scope: %s\n", result.TrustScope)
	for _, condition := range result.ForensicConditions {
		explainForensicCondition(condition, stdout)
	}
	explainReasons(result.Reasons, stdout)
	explainNextActions(result.NextActions, stdout)
	return 0
}
