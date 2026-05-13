package demo

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func requiredEvidenceIDs(contract trace.Contract) []string {
	// requiredEvidenceIDs keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	ids := make([]string, 0, len(contract.RequiredEvidence))
	for _, requirement := range contract.RequiredEvidence {
		if requirement.ID != "" {
			ids = append(ids, requirement.ID)
		}
	}
	return ids
}
