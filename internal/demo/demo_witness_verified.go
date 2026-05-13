package demo

import (
	"sort"
)

func ciWitnessVerified(record WitnessSummary) bool {

	return record.Kind == "github-actions" && record.Status == GatePass && record.TrustScope == "ci_witnessed"
}

func applyVerifiedCIWitness(result GateResult, record WitnessSummary, expected WitnessExpectation) GateResult {
	// applyVerifiedCIWitness keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	if bindingState, bindingReasons := witnessBindingState(record, expected); bindingState != GatePass {
		result.CIWitnessGate = bindingState
		result.Reasons = append(result.Reasons, bindingReasons...)
		for _, reason := range bindingReasons {
			result.WitnessBindings = append(result.WitnessBindings, WitnessBinding{ID: "source", State: bindingState, Reason: reason})
		}
		sort.Strings(result.Reasons)
		result.GateConditions = gateConditions(result)
		return result
	}

	result.CIWitnessGate = GatePass
	result.MissingAuditEvidence = []string{"external_witness_checkpoint"}
	result.GateConditions = gateConditions(result)
	return result
}
