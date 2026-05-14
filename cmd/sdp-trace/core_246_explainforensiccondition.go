package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/forensic"
)

func explainForensicCondition(condition forensic.Condition, stdout io.Writer) {
	fmt.Fprintf(stdout, "Forensic condition %s: %s (%s)\n", condition.ID, condition.State, condition.ReasonCode)
	if condition.CappedToRetentionMode != "" {
		// The cap belongs to this condition, not to the entire artifact.
		fmt.Fprintf(stdout, "Capped to retention mode: %s\n", condition.CappedToRetentionMode)
	}
}
