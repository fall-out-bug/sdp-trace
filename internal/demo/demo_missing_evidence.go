package demo

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func missingContractEvidence(rows []RunRow, contract trace.Contract) []string {
	// missingContractEvidence keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	observed := observedEvidenceKinds(rows)
	missing := make([]string, 0)
	for _, requirement := range contract.RequiredEvidence {
		if requirement.ID == "" || observed[requirement.ID] {
			continue
		}
		missing = append(missing, requirement.ID)
	}
	return missing
}

func observedEvidenceKinds(rows []RunRow) map[string]bool {
	// observedEvidenceKinds keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	observed := map[string]bool{}
	for _, row := range rows {
		if rowHasObservedEvidenceKind(row) {
			observed[row.Kind] = true
		}
	}
	return observed
}

func rowHasObservedEvidenceKind(row RunRow) bool {

	return row.Kind != "" && row.Kind != "unmatched" && row.Result == trace.VerdictObserved
}
