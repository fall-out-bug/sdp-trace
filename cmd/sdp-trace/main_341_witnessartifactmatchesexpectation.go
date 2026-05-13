package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func witnessArtifactMatchesExpectation(artifact demo.WitnessArtifactDigest, expectedArtifacts map[string]string) bool {
	expectedSHA, ok := expectedArtifacts[artifact.Path]
	return ok && expectedSHA == artifact.SHA256
}
