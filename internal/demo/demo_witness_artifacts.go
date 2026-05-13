package demo

import (
	"fmt"
)

func witnessArtifactBindingState(actual, expected []WitnessArtifactDigest) (string, []string) {
	// witnessArtifactBindingState keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	expectedArtifacts := witnessArtifactsByPath(expected)
	seenArtifacts := map[string]bool{}
	for _, artifact := range actual {
		seenArtifacts[artifact.Path] = true
		if state, reasons := validateWitnessArtifact(artifact, expectedArtifacts); state != GatePass {
			return state, reasons
		}
	}
	return missingWitnessArtifactState(expectedArtifacts, seenArtifacts)
}

func missingWitnessArtifactState(expectedArtifacts map[string]string, seenArtifacts map[string]bool) (string, []string) {
	// missingWitnessArtifactState keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	for path := range expectedArtifacts {
		if !seenArtifacts[path] {
			return GateCannotVerify, []string{fmt.Sprintf("ci witness artifact %s is missing from witness", path)}
		}
	}
	return GatePass, nil
}

func witnessArtifactsByPath(artifacts []WitnessArtifactDigest) map[string]string {
	// witnessArtifactsByPath keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	byPath := map[string]string{}
	for _, artifact := range artifacts {
		byPath[artifact.Path] = artifact.SHA256
	}
	return byPath
}

func validateWitnessArtifact(artifact WitnessArtifactDigest, expected map[string]string) (string, []string) {
	// validateWitnessArtifact keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	expectedDigest, ok := expected[artifact.Path]
	if !ok {
		return GateCannotVerify, []string{fmt.Sprintf("ci witness artifact %s is not present in current gate input", artifact.Path)}
	}
	if expectedDigest != artifact.SHA256 {
		return GateFail, []string{fmt.Sprintf("ci witness artifact digest mismatch for %s", artifact.Path)}
	}
	return GatePass, nil
}
