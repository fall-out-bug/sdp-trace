package ciartifact

func evaluateIndex(input ArtifactIndexInput) ArtifactIndexResult {
	// Index evaluation checks whether the manifest index itself is present and safe
	// before any artifact-family evidence is trusted.

	state := defaultString(input.State, IndexNotAssessed)
	if outcome, ok := indexOutcomes[state]; ok {
		return ArtifactIndexResult{State: state, Result: outcome.state, ReasonCode: outcome.code, Reason: outcome.reason}
	}
	outcome := indexOutcomes[IndexUnverifiable]

	return ArtifactIndexResult{State: IndexUnverifiable, Result: outcome.state, ReasonCode: outcome.code, Reason: "artifact index state is unrecognized under selected profile"}
}

var indexOutcomes = map[string]familyOutcome{
	IndexValid:          {StatePass, "artifact_index_valid", "artifact index is present and valid", ""},
	IndexSelfReference:  {StateFail, "artifact_index_self_reference", "artifact index includes itself as an indexed entry", ""},
	IndexDigestMismatch: {StateFail, "artifact_digest_mismatch", "artifact digest contradicts selected artifact metadata", ""},
	IndexMissing:        {StateCannotVerify, "artifact_index_missing", "required artifact index is missing", ""},
	IndexUnverifiable:   {StateCannotVerify, "artifact_index_unverifiable", "artifact index could not be verified", ""},
	IndexNotAssessed:    {StateNotAssessed, "artifact_index_not_assessed", "artifact index was outside selected profile scope", ""},
}
