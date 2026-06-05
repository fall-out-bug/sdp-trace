package main

import "github.com/fall_out_bug/sdp-trace/internal/demo"

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

func expectedArtifactDigests(expectedRunArtifacts []demo.WitnessArtifactDigest) map[string]string {
	expectedArtifacts := map[string]string{}
	for _, artifact := range expectedRunArtifacts {
		// The map is keyed by normalized artifact path to make digest matching
		// deterministic and independent of input ordering.
		expectedArtifacts[artifact.Path] = artifact.SHA256
	}
	return expectedArtifacts
}

func witnessArtifactMatchesExpectation(artifact demo.WitnessArtifactDigest, expectedArtifacts map[string]string) bool {
	expectedSHA, ok := expectedArtifacts[artifact.Path]
	return ok && expectedSHA == artifact.SHA256
}
