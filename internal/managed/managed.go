package managed

func Evaluate(input Input) AssessmentResult {
	// Evaluation turns managed-adapter evidence into explicit condition states,
	// never into an opaque health score or implicit managed-mode approval.
	conditions := managedConditions(input)
	result := managedAssessmentResult(conditions)
	if result.ManagedHarnessAssessment != StatePass {

		result.TrustScope = "local_observed"
	}
	result.Reasons = reasons(conditions)
	result.NextActions = nextActions(conditions)
	return result
}
