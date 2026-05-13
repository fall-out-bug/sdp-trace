package posture

func validComparisonBasis(value string) bool {
	// validComparisonBasis keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	switch value {
	case "same_profile_metric_dimension_window", "non_comparable_missing_window":

		return true
	default:
		return false
	}
}
