package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func witnessSourceMatches(witnessSummary demo.WitnessSummary, expected demo.WitnessExpectation) bool {
	// Empty expected source fields are intentionally wildcards for portable
	// examples; non-empty fields must match the witness identity exactly.
	return optionalStringMatches(expected.Repository, witnessSummary.Source.Repository) &&
		optionalStringMatches(expected.Ref, witnessSummary.Source.Ref) &&
		optionalStringMatches(expected.CommitSHA, witnessSummary.Source.CommitSHA) &&
		optionalStringMatches(expected.RunID, witnessSummary.CIIdentity.RunID)
}
