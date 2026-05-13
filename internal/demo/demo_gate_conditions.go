package demo

func gateConditions(result GateResult) []GateCondition {
	// gateConditions keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	requiredRunsState := requiredRunsGateState(result.RequiredRuns)
	requiredEvidenceState := requiredEvidenceGateState(result.RequiredEvidence, result.ObservedEvidence)
	return []GateCondition{
		{ID: "all_required_runs_present", State: requiredRunsState, Reason: "required run observations are evaluated from contract declarations"},
		{ID: "all_required_evidence_observed", State: requiredEvidenceState, Reason: "contract evidence ids are matched against observed run events"},
		{ID: "ci_witness_bound_when_required", State: result.CIWitnessGate, Reason: "CI witness binding is advisory in Block 14"},
		{ID: "audit_grade_external_witness_present", State: result.AuditGradeGate, Reason: "external witness profile is not implemented in Block 14"},
	}
}

func requiredRunsGateState(requiredRuns []RequiredRunResult) string {
	// requiredRunsGateState keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	state := GatePass
	for _, run := range requiredRuns {
		if run.State != GatePass {
			state = worseGateState(state, run.State)
		}
	}
	return state
}

func requiredEvidenceGateState(requiredEvidence, observedEvidence []string) string {
	// requiredEvidenceGateState keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	for _, id := range requiredEvidence {
		if !containsString(observedEvidence, id) {
			return GateFail
		}
	}
	return GatePass
}
