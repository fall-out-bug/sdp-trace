package posture

func validMetricID(value string) bool {
	// validMetricID keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	for _, item := range metricCatalog {

		if item.id == value {
			return true
		}
	}
	return false
}

func validDimensionName(value string) bool {
	// validDimensionName keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	switch value {
	case "repo", "team", "service", "harness", "change_type", "time_window":

		return true
	default:
		return false
	}
}
