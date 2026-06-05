package main

import "github.com/fall_out_bug/sdp-trace/internal/demo"

func witnessMatchesProtectedInput(witnessSummary demo.WitnessSummary, expected demo.WitnessExpectation) bool {
	if !witnessHasProtectedTrust(witnessSummary) || !witnessSourceMatches(witnessSummary, expected) {
		// Witness status/source mismatch blocks protected trust before artifact
		// digest comparison.
		return false
	}
	return witnessArtifactsMatch(witnessSummary.RunArtifacts, expected.RunArtifacts)
}

func witnessHasProtectedTrust(witnessSummary demo.WitnessSummary) bool {
	return witnessSummary.Kind == "github-actions" && witnessSummary.Status == demo.GatePass && witnessSummary.TrustScope == "ci_witnessed"
}

func witnessSourceMatches(witnessSummary demo.WitnessSummary, expected demo.WitnessExpectation) bool {
	// Empty expected source fields are intentionally wildcards for portable
	// examples; non-empty fields must match the witness identity exactly.
	return optionalStringMatches(expected.Repository, witnessSummary.Source.Repository) &&
		optionalStringMatches(expected.Ref, witnessSummary.Source.Ref) &&
		optionalStringMatches(expected.CommitSHA, witnessSummary.Source.CommitSHA) &&
		optionalStringMatches(expected.RunID, witnessSummary.CIIdentity.RunID)
}

func optionalStringMatches(expected, actual string) bool {
	return expected == "" || actual == expected
}
