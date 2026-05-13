package ciartifact

func lowerAuthorityReason(producer string) string {
	// lowerAuthorityReason keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if reason, ok := lowerAuthorityReasons[producer]; ok {
		return reason
	}

	return "lower_authority_producer_scope"
}

var lowerAuthorityReasons = map[string]string{
	ProducerCheckedIn:           "checked_in_claim_contradicts_ci_artifacts",
	ProducerAgentReported:       "agent_reported_claim_without_observed_family",
	ProducerExternalArtifactRef: "external_artifact_ref_unverifiable",
}
