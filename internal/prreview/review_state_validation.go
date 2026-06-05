package prreview

// Review state allowlists stay narrow because they feed trust decisions.
//
// Unknown CI states or reviewer runners must fail validation instead of being
// normalized into a green or cannot-verify state.
func validCIState(state string) bool {
	switch state {
	case StatePass, StateFail, StatePending, StateNotAssessed, StateCannotVerify:
		return true
	default:
		return false
	}
}

func validRunner(runner string) bool {
	switch runner {
	case RunnerPI, RunnerOpenCode, RunnerManualExternal:
		return true
	default:
		return false
	}
}
