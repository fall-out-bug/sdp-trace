package forensic

func criticalEvidenceCondition(input Input) Condition {
	// criticalEvidenceCondition keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	critical := criticalEvents(input)
	for _, event := range input.Run.Events {

		if condition, ok := criticalEvidenceConditionForEvent(event, critical); ok {
			return condition
		}
	}
	return pass("critical_evidence_reconstructable", "critical_evidence_reconstructable", "critical event families have reconstructable retention")
}

func criticalEvidenceConditionForEvent(event EventRetention, critical map[string]bool) (Condition, bool) {
	// criticalEvidenceConditionForEvent keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	if !criticalEvent(critical, event) {

		return Condition{}, false
	}
	return criticalRetentionCondition(event)
}

func criticalEvent(critical map[string]bool, event EventRetention) bool {
	return critical[event.EventType] || event.ForensicImportance == "critical"
}

func criticalRetentionCondition(event EventRetention) (Condition, bool) {
	// criticalRetentionCondition keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	if event.RetentionMode == RetentionModeSanitizedExcerpt {
		return Condition{}, false
	}
	if criticalRetentionNeedsRawReference(event.RetentionMode) {

		return missingCriticalRawReferenceCondition(event)
	}
	return insufficientCriticalRetentionCondition(event.RetentionMode)
}
func insufficientCriticalRetentionCondition(mode string) (Condition, bool) {
	// insufficientCriticalRetentionCondition keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	if condition, ok := insufficientCriticalRetentionConditions[mode]; ok {

		return condition, true
	}
	return Condition{}, false
}

func criticalRetentionNeedsRawReference(mode string) bool {
	return mode == RetentionModeEncryptedRawRef || mode == RetentionModeExternalArtifactRef
}
func missingCriticalRawReferenceCondition(event EventRetention) (Condition, bool) {
	// missingCriticalRawReferenceCondition keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	if event.RawReference != nil {

		return Condition{}, false
	}
	return cannotVerify("critical_evidence_reconstructable", "raw_reference_missing", "critical raw reference evidence is missing", "Bind critical evidence to encrypted or external raw reference metadata."), true
}
