package demo

import (
	"fmt"
)

func applyWitness(result GateResult, witnessPath string) GateResult {

	return applyWitnessWithExpectation(result, witnessPath, WitnessExpectation{})
}
func applyWitnessWithExpectation(result GateResult, witnessPath string, expected WitnessExpectation) GateResult {
	// applyWitnessWithExpectation keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	record, err := loadWitnessSummary(witnessPath)
	if err != nil {
		result.CIWitnessGate = GateCannotVerify
		result.Reasons = append(result.Reasons, fmt.Sprintf("ci witness cannot verify: %v", err))
		return result
	}
	result.Witness = &record
	if ciWitnessVerified(record) {

		return applyVerifiedCIWitness(result, record, expected)
	}
	result.CIWitnessGate = GateCannotVerify
	result.MissingAuditEvidence = []string{"ci_oidc_witness", "external_witness_checkpoint"}
	result.GateConditions = gateConditions(result)
	return result
}
