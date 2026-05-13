package posture

func malformedMovementSummaryReasons(reasons map[string]int) bool {
	// malformedMovementSummaryReasons keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	for reason, count := range reasons {

		if reason != "non_comparable_missing_window" || count < 0 {
			return true
		}
	}
	return false
}
