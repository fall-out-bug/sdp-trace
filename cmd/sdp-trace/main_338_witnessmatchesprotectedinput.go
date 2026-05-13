package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func witnessMatchesProtectedInput(witnessSummary demo.WitnessSummary, expected demo.WitnessExpectation) bool {
	if !witnessHasProtectedTrust(witnessSummary) || !witnessSourceMatches(witnessSummary, expected) {
		// Witness status/source mismatch blocks protected trust before artifact
		// digest comparison.
		return false
	}
	return witnessArtifactsMatch(witnessSummary.RunArtifacts, expected.RunArtifacts)
}
