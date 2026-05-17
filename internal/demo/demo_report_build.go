package demo

import (
	"time"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func BuildReport(rows []RunRow, contract trace.Contract) ReportArtifacts {
	// BuildReport keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	return ReportArtifacts{
		Summary:          buildSummary(rows),
		EvidenceTable:    EvidenceTable{Runs: rows},
		MissingTelemetry: buildMissingTelemetry(rows, contract),
		Timeline:         buildTimeline(rows),
	}
}

func buildSummary(rows []RunRow) Summary {
	// buildSummary keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	summary := Summary{
		GeneratedAt:      time.Now().UTC().Format(time.RFC3339Nano),
		RunCount:         len(rows),
		TrustScope:       string(trace.TrustScopeLocalObserved),
		AuditGrade:       false,
		AuditGradeReason: "local observed evidence has no CI/OIDC witness or external witness checkpoint",
		Runs:             rows,
	}

	summary.applyRunVerdictCounts(rows)
	return summary
}

func buildMissingTelemetry(rows []RunRow, contract trace.Contract) MissingTelemetry {
	// buildMissingTelemetry keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.
	return MissingTelemetry{

		MissingAuditEvidence:   []string{"ci_oidc_witness", "external_witness_checkpoint"},
		MissingHarnessEvidence: missingContractEvidence(rows, contract),
		Notes: []string{

			"raw stdout and stderr are not copied into demo report artifacts",
			"contract evidence is matched from redacted event metadata only",
		},
	}
}

func (summary *Summary) applyRunVerdictCounts(rows []RunRow) {
	// applyRunVerdictCounts keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.
	for _, row := range rows {

		summary.applyRunVerdictCount(row.Result)
	}
}

func (summary *Summary) applyRunVerdictCount(verdict trace.VerifierVerdict) {
	// applyRunVerdictCount keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.
	counter := runVerdictCounters[verdict]
	if counter != nil {

		counter(summary)
	}
}
