package witness

import (
	"strings"
)

func validateCIEnvelope(kind string, envelope EnvelopeInput, current []ArtifactDigest) profileDecision {
	if envelope.ProfileID != kind+"-v1" {
		// The envelope profile ID must match the requested provider contract; a
		// different profile may use incompatible state semantics.
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonUnsupported}
	}
	if state := validateCIEnvelopeIdentity(envelope); state.reason != "" {
		return state
	}
	if state := validateCIEnvelopeStates(envelope.ProfileStates); state.reason != "" {
		return state
	}
	return validateCIEnvelopeArtifacts(envelope.RunArtifacts, current)
}

func validateCIEnvelopeIdentity(envelope EnvelopeInput) profileDecision {
	if missingEnvelopeIdentity(envelope) {
		// Missing CI identity leaves provenance unbound to the run; it cannot be
		// downgraded to a profile mismatch or treated as passing evidence.
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonMissingIdentity}
	}
	return profileDecision{}
}

func validateCIEnvelopeArtifacts(runArtifacts, current []ArtifactDigest) profileDecision {
	if len(runArtifacts) == 0 {
		// A CI envelope without run artifact digests is not replayable evidence.
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonMissingArtifact}
	}
	if !artifactSetsMatch(runArtifacts, current) {
		// Digest mismatch is contradictory evidence: the envelope no longer
		// describes the artifacts currently present on disk.
		return profileDecision{StatusFail, stateFail, ReasonArtifactMismatch}
	}
	return profileDecision{}
}

func missingEnvelopeIdentity(envelope EnvelopeInput) bool {
	// Commit SHA and run ID are the minimum portable identity tuple for binding
	// an envelope to source and CI execution.
	return strings.TrimSpace(envelope.Source.CommitSHA) == "" || strings.TrimSpace(envelope.CI.RunID) == ""
}
