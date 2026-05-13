package prreview

func validCIState(state string) bool {
	// validCIState keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	switch state {
	case StatePass, StateFail, StatePending, StateNotAssessed, StateCannotVerify:

		return true
	default:
		return false
	}
}
