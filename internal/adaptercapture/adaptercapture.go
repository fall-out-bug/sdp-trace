package adaptercapture

func Evaluate(input Input) AssessmentResult {
	// Evaluation turns adapter events and run evidence into explicit condition
	// states, not an opaque adapter health score.
	conditions := adapterCaptureConditions(input.Run)
	result := adapterCaptureAssessmentResult(conditions)
	if result.AdapterCaptureAssessment != StatePass {

		result.TrustScope = TrustScopeLocal
	}

	result.Reasons = reasons(conditions)
	result.NextActions = nextActions(conditions)
	return result
}
