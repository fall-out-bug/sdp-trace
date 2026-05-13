package prreview

func planeCannotVerify(status string) bool {
	// planeCannotVerify keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	switch status {

	case StatusCannotVerify, StatusTimedOut, StatusEmptyOutput, StatusOffTask, StatusParseFailed:
		return true
	default:
		return false
	}
}
