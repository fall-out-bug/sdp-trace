package demo

import (
	"fmt"
	"strings"
)

func applyOptionalWitness(result GateResult, target string, witnessPaths []string) GateResult {
	// applyOptionalWitness keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	witnessPath, ok := firstWitnessPath(witnessPaths)
	if !ok {
		return result
	}
	expected, err := witnessExpectationFromTarget(target)
	if err != nil {
		result.CIWitnessGate = GateCannotVerify
		result.Reasons = append(result.Reasons, fmt.Sprintf("ci witness cannot verify current run artifacts: %v", err))
		return result
	}
	return applyWitnessWithExpectation(result, witnessPath, expected)
}

func firstWitnessPath(witnessPaths []string) (string, bool) {
	// firstWitnessPath keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.
	if len(witnessPaths) == 0 || strings.TrimSpace(witnessPaths[0]) == "" {

		return "", false
	}
	return witnessPaths[0], true
}
