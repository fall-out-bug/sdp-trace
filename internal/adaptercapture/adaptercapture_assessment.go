package adaptercapture

func adapterCaptureConditions(run RunEvidence) []Condition {
	// Conditions stay grouped by contract, identity, run binding, task, tool depth,
	// mutation, model identity, test provenance, redaction, and overclaim evidence.

	return []Condition{
		contractCondition(run),
		identityCondition(run),
		runBindingCondition(run),
		taskDriftCondition(run),
		toolDepthCondition(run),
		fileMutationCondition(run),
		modelIdentityCondition(run),
		testProvenanceCondition(run),
		providerRefsCondition(run),
		redactionMetadataCondition(run),
		overclaimCondition(run),
	}
}

func adapterCaptureAssessmentResult(conditions []Condition) AssessmentResult {
	// Result assembly keeps top-level state, reasons, and actions derived from
	// machine conditions rather than prose-only claims.

	return AssessmentResult{
		SchemaVersion:            SchemaVersion,
		SelectedProfile:          ProfileAdapterCapture,
		AdapterCaptureAssessment: topLevel(conditions),
		TrustScope:               TrustScopeAdapter,
		AdapterCaptureConditions: conditions,
	}
}
