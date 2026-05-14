package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/adaptercapture"
)

func explainAdapterCaptureCondition(condition adaptercapture.Condition, stdout io.Writer) {
	fmt.Fprintf(stdout, "Adapter condition %s: %s (%s)\n", condition.ID, condition.State, condition.ReasonCode)
	if condition.CappedToRetentionMode != "" {
		// Retention caps are printed beside the condition that caused them so
		// the explanation preserves the evidence chain.
		fmt.Fprintf(stdout, "Capped to retention mode: %s\n", condition.CappedToRetentionMode)
	}
}
