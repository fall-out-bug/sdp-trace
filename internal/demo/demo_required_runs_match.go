package demo

import (
	"fmt"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func matchRequiredRun(row RunRow, required trace.RequiredRun, result RequiredRunResult) RequiredRunResult {
	// matchRequiredRun keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	result.MatchedRunID = row.RunID
	result.State = GatePass
	result.Reasons = []string{fmt.Sprintf("required run %s matched wrapper %s", required.ID, required.WrapperName)}
	if row.Result != trace.VerdictObserved || row.ClosureState != trace.ClosureStateCompleted {
		result = cannotVerifyRequiredRun(result, required.ID, row.Name)
	}
	if evidenceID, ok := missingEvidenceID(row, required.RequiredEvidence); ok {
		return cannotVerifyRequiredRunEvidence(result, required.ID, evidenceID)
	}
	return result
}

func missingEvidenceID(row RunRow, requiredEvidence []string) (string, bool) {
	// missingEvidenceID keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	for _, evidenceID := range requiredEvidence {
		if row.Kind != evidenceID {
			return evidenceID, true
		}
	}
	return "", false
}
