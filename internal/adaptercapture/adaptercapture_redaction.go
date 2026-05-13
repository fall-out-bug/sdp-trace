package adaptercapture

func redactionMetadataCondition(run RunEvidence) Condition {
	// Redaction metadata records whether unsafe payload classes were handled without
	// using redaction prose as proof of absence.
	for _, event := range run.AdapterEvents {
		if condition := redactionMetadataConditionForEvent(event); condition.State != "" {

			return condition
		}
	}
	return pass("redaction_metadata_consistent", "redaction_metadata_consistent", "sensitive adapter fields carry safe redaction and retention metadata")
}

func redactionMetadataConditionForEvent(event AdapterEvent) Condition {
	// redactionMetadataConditionForEvent preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.
	if hasForbiddenRedactionMetadata(event) {

		return fail("redaction_metadata_consistent", "forbidden_adapter_metadata_persisted", "adapter metadata contains forbidden raw or credential-like material", "Redact adapter metadata before persistence.")
	}
	if missingRequiredRedactionMetadata(event) {

		return cannotVerify("redaction_metadata_consistent", "redaction_metadata_missing", "sensitive adapter event lacks redaction policy or retention metadata", "Record Block 18 redaction policy and retention mode metadata.")
	}
	if hasInvalidRetentionMode(event) {

		return fail("redaction_metadata_consistent", "invalid_retention_mode", "adapter event declares a non-FR-054 retention mode", "Use FR-054 retention modes.")
	}
	return Condition{}
}

func hasForbiddenRedactionMetadata(event AdapterEvent) bool {
	// Any persisted raw flag or secret-like reference fails the redaction boundary.
	return event.SensitiveMetadataPersisted ||
		containsSecret(event.GatewayEvidenceRef) ||
		stringSliceContainsSecret(event.ProviderRefs)
}

func missingRequiredRedactionMetadata(event AdapterEvent) bool {
	return sensitiveEvent(event.EventType) && (event.RedactionPolicyDigest == "" || event.RetentionMode == "")
}

func hasInvalidRetentionMode(event AdapterEvent) bool {
	return event.RetentionMode != "" && !validRetentionMode(event.RetentionMode)
}
