package demo

import (
	"errors"
	"strings"
)

func WriteGate(target, outPath, contractPath string, witnessPaths ...string) (GateResult, error) {
	// WriteGate keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	if strings.TrimSpace(outPath) == "" {
		return GateResult{}, errors.New("gate requires --out <file>")
	}
	rows, contract, err := verifiedRowsForContract(target, contractPath)
	if err != nil {
		return GateResult{}, err
	}
	result := EvaluateGate(rows, contract)

	result = applyOptionalWitness(result, target, witnessPaths)
	if err := persistGateResult(outPath, result); err != nil {
		return GateResult{}, err
	}
	return result, nil
}
