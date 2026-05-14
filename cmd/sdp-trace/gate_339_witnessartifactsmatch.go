package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func witnessArtifactsMatch(runArtifacts, expectedRunArtifacts []demo.WitnessArtifactDigest) bool {
	// Artifact counts must match exactly so extra or missing witness artifacts
	// cannot hide under a partial digest match.
	expectedArtifacts := expectedArtifactDigests(expectedRunArtifacts)
	if len(runArtifacts) != len(expectedArtifacts) {
		return false
	}
	for _, artifact := range runArtifacts {
		// Every reported witness artifact must match the expected path and
		// digest derived from the observed run.
		if !witnessArtifactMatchesExpectation(artifact, expectedArtifacts) {
			return false
		}
	}
	return true
}
