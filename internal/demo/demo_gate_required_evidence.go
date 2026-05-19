package demo

import (
	"fmt"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func applyRequiredEvidence(result *GateResult, contract trace.Contract, observedEvidence map[string]bool) {
	// applyRequiredEvidence keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	for _, requirement := range contract.RequiredEvidence {
		if observedEvidence[requirement.ID] {
			result.ObservedEvidence = append(result.ObservedEvidence, requirement.ID)
			continue
		}
		result.LocalGate = GateFail
		result.Reasons = append(result.Reasons, fmt.Sprintf("missing locally observed contract evidence %s", requirement.ID))
	}
}
