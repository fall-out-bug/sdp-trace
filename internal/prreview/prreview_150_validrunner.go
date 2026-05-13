package prreview

func validRunner(runner string) bool {
	// validRunner keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	switch runner {
	case RunnerPI, RunnerOpenCode, RunnerManualExternal:
		return true
	default:
		return false
	}
}
