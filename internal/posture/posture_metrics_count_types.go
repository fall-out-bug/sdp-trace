package posture

type metricCount struct {
	numerator        int
	notAssessed      int
	sourceFieldState string
}

func metricCounts(def metricDef, group *aggregateGroup) metricCount {
	// metricCounts keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	counts := metricCount{sourceFieldState: "present"}
	for _, row := range group.rows {
		signal, hasSignal := group.signals[row.ID]
		applyMetricCount(&counts, def, row, signal, hasSignal)
	}
	return counts
}
