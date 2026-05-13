package demo

import (
	"strings"
)

func protectedWitnessFreshnessCondition(input ProtectedGateInput) ProtectedCondition {
	// protectedWitnessFreshnessCondition keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	generatedAt, ok := protectedWitnessGeneratedAt(input.Witness)
	if !ok {
		return witnessFreshnessCannotVerify("missing_witness_freshness", "CI witness generated_at is required for protected freshness evaluation", "Supply CI witness evidence with generated_at freshness data.")
	}
	return protectedWitnessFreshnessAt(generatedAt, input.Now)
}
func protectedWitnessGeneratedAt(witness *WitnessSummary) (string, bool) {
	// protectedWitnessGeneratedAt keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.
	if witness == nil || strings.TrimSpace(witness.GeneratedAt) == "" {

		return "", false
	}
	return witness.GeneratedAt, true
}
