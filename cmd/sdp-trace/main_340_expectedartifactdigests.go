package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func expectedArtifactDigests(expectedRunArtifacts []demo.WitnessArtifactDigest) map[string]string {
	expectedArtifacts := map[string]string{}
	for _, artifact := range expectedRunArtifacts {
		// The map is keyed by normalized artifact path to make digest matching
		// deterministic and independent of input ordering.
		expectedArtifacts[artifact.Path] = artifact.SHA256
	}
	return expectedArtifacts
}
