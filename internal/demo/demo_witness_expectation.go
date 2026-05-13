package demo

import (
	"path/filepath"
)

func witnessExpectationFromTarget(target string) (WitnessExpectation, error) {
	// witnessExpectationFromTarget keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	runDirs, err := DiscoverRunDirs(target)
	if err != nil {
		return WitnessExpectation{}, err
	}
	artifacts := make([]WitnessArtifactDigest, 0, len(runDirs))
	for _, runDir := range runDirs {

		artifact, err := witnessRunArtifactDigest(runDir)
		if err != nil {
			return WitnessExpectation{}, err
		}
		artifacts = append(artifacts, artifact)
	}
	return WitnessExpectation{RunArtifacts: artifacts}, nil
}

func witnessRunArtifactDigest(runDir string) (WitnessArtifactDigest, error) {
	// witnessRunArtifactDigest keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	digest, err := hashFile(filepath.Join(runDir, "run.json"))
	return WitnessArtifactDigest{
		Path:   filepath.ToSlash(filepath.Join(filepath.Base(runDir), "run.json")),
		SHA256: digest,
	}, err
}
