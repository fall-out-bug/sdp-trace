package witness

import (
	"strings"
)

func unsafeEnvelopeFields(envelope EnvelopeInput) bool {
	// Envelope safety covers both identity metadata and artifact references
	// because either can leak credentials or private paths.
	return unsafeEnvelopeScalarFields(envelope) || unsafeEnvelopeArtifactFields(envelope)
}

func unsafeEnvelopeScalarFields(envelope EnvelopeInput) bool {
	// Scalar provider fields are copied into witness records, so they must be
	// safe before the envelope can influence the output record.
	// Structured artifact references are checked separately because their path
	// and digest rules differ.
	// The list mirrors fields copied by applyCIEnvelopeRecordValues.
	// Any unsafe scalar blocks the entire envelope before profile states are
	// trusted.
	values := []string{
		envelope.Source.Repository,
		envelope.Source.Ref,
		envelope.Source.CommitSHA,
		envelope.CI.Provider,
		envelope.CI.ServerURL,
		envelope.CI.Workflow,
		envelope.CI.Job,
		envelope.CI.RunID,
		envelope.CI.RunAttempt,
		envelope.CI.Actor,
	}
	for _, value := range values {
		if unsafeOutputString(value) {
			return true
		}
	}
	return false
}

func unsafeEnvelopeArtifactFields(envelope EnvelopeInput) bool {
	for _, artifact := range append(envelope.RunArtifacts, envelope.ReportArtifacts...) {
		// Artifact paths and digests are persisted verbatim in the witness output.
		if unsafeOutputString(artifact.Path) || unsafeOutputString(artifact.SHA256) {
			return true
		}
	}
	return false
}

func unsafeOutputString(value string) bool {
	if value == "" {
		return false
	}
	if containsSecretLike([]byte(value)) {
		// Secret-shaped scalar values fail before broader personal/path checks.
		return true
	}
	lower := strings.ToLower(value)
	return strings.Contains(lower, "/private/") || strings.Contains(lower, "@")
}
