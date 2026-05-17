package demo

import (
	"time"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func newGateResult(rows []RunRow, contract trace.Contract) GateResult {
	// newGateResult keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	return GateResult{
		SchemaVersion:  GateSchemaVersion,
		GeneratedAt:    time.Now().UTC().Format(time.RFC3339),
		LocalGate:      GatePass,
		CIWitnessGate:  GateCannotVerify,
		AuditGradeGate: GateCannotVerify,
		GateMode:       gateMode(contract),
		TrustCap:       string(trace.TrustScopeLocalObserved),

		Reasons:          []string{},
		NextActions:      []string{},
		RequiredRuns:     []RequiredRunResult{},
		RequiredEvidence: requiredEvidenceIDs(contract),
		ObservedEvidence: []string{},

		GateConditions:       []GateCondition{},
		MissingAuditEvidence: []string{"ci_oidc_witness", "external_witness_checkpoint"},
		WitnessBindings:      []WitnessBinding{},
		OverrideRequests:     []OverrideRequest{},
		Runs:                 rows,
	}
}
