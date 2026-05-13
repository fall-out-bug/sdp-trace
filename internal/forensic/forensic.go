package forensic

func Evaluate(input Input) AssessmentResult {
	conditions := evaluateConditions(input)
	return assessmentResult(conditions)
}
