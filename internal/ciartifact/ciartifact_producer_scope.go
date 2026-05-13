package ciartifact

import "strings"

func safeRequiredProducerScope(value string) string {
	// safeRequiredProducerScope keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if strings.TrimSpace(value) == "" {

		return ProducerCIUploaded
	}
	scope := safeProducerScope(value)
	if scope == ProducerNotAssessed && value != ProducerNotAssessed {

		return ProducerCIUploaded
	}
	return scope
}

func safeAuthorityScope(value string) string {
	// safeAuthorityScope keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if safeToken(value) {
		return defaultString(value, AuthorityScopeObservation)
	}

	return AuthorityScopeObservation
}

func safeToken(value string) bool {
	return len(value) > 0 && len(value) <= 128 && safeIdentityToken(value, "_.:-")
}
