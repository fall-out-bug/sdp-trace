package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/adaptercapture"
)

func explainAdapterCaptureAssessment(result adaptercapture.AssessmentResult, stdout io.Writer) int {
	// Adapter explanations are derived views of stored conditions, reasons,
	// and next actions; the assessment verdict is not recomputed here.
	fmt.Fprintf(stdout, "Selected profile: %s\n", result.SelectedProfile)
	fmt.Fprintf(stdout, "Adapter capture assessment: %s\n", result.AdapterCaptureAssessment)
	fmt.Fprintf(stdout, "Trust scope: %s\n", result.TrustScope)
	for _, condition := range result.AdapterCaptureConditions {
		explainAdapterCaptureCondition(condition, stdout)
	}
	explainReasons(result.Reasons, stdout)
	explainNextActions(result.NextActions, stdout)
	return 0
}
