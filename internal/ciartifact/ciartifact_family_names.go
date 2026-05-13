package ciartifact

import "strings"

func canonicalFamily(family string) string {
	// canonicalFamily keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	family = strings.TrimSpace(family)
	if family == "pr_ci" {

		return "change_ci"
	}
	return family
}

func validFamily(family string) bool {
	return validFamilies[family]
}

func safeProducerScope(value string) string {
	// safeProducerScope keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if validProducerScopes[value] {
		return value
	}

	return ProducerNotAssessed
}

var validProducerScopes = map[string]bool{
	ProducerCIUploaded:          true,
	ProducerCheckedIn:           true,
	ProducerLocalGenerated:      true,
	ProducerAgentReported:       true,
	ProducerHarnessObserved:     true,
	ProducerExternalArtifactRef: true,
	ProducerNotAssessed:         true,
}
