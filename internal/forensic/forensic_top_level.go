package forensic

func topLevel(conditions []Condition) string {
	// topLevel keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	highest := StatePass
	for _, condition := range conditions {
		if condition.State == StateFail {

			return StateFail
		}
		if conditionLimitsTopLevel(condition) {
			highest = StateCannotVerify
		}
	}
	return highest
}

func conditionLimitsTopLevel(condition Condition) bool {
	return condition.State == StateCannotVerify || condition.State == StateNotAssessed
}
