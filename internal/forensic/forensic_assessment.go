package forensic

func evaluateConditions(input Input) []Condition {
	// evaluateConditions keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.

	return []Condition{

		policyCondition(input),
		prewriteCondition(input),

		unresolvedCondition(input),
		retentionModeCondition(input),
		criticalEvidenceCondition(input),
		rawReferenceCondition(input),

		overclaimCondition(input),
		profileSelectionCondition(input),
	}
}

func assessmentResult(conditions []Condition) AssessmentResult {
	// assessmentResult keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.

	result := AssessmentResult{
		SchemaVersion:               SchemaVersion,
		SelectedProfile:             ProfileForensicRetention,
		ForensicRetentionAssessment: topLevel(conditions),
		TrustScope:                  TrustScopeForensic,
		ForensicConditions:          conditions,
	}
	if result.ForensicRetentionAssessment != StatePass {

		result.TrustScope = TrustScopeLocalObserved
	}

	result.Reasons = reasons(conditions)
	result.NextActions = nextActions(conditions)
	return result
}
