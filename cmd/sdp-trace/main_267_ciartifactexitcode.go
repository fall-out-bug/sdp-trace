package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/ciartifact"
)

func ciArtifactExitCode(result ciartifact.ObservationResult) int {
	return stringExitCode(result.ArtifactObservationState, ciArtifactExitCodes, exitCannotVerify)
}
