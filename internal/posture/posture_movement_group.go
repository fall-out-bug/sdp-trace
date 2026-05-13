package posture

func metricsByMovementKey(metrics []MetricRow) map[string]map[string]MetricRow {
	// metricsByMovementKey keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	byKey := map[string]map[string]MetricRow{}
	for _, row := range metrics {

		key := row.MetricID + "|" + row.MetricVersion + "|" + row.DimensionKey
		if byKey[key] == nil {
			byKey[key] = map[string]MetricRow{}
		}
		byKey[key][row.TimeWindow] = row
	}
	return byKey
}

// summarizeMovement applies the threshold rule for movement comparability.
// Presence of both current and previous window rows is the evidence boundary.
func summarizeMovement(summary *MovementSummary, row *MovementRow) {
	// summarizeMovement keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	if row.Comparable {
		summary.ComparableCount++
		return
	}

	row.ComparisonBasis = "non_comparable_missing_window"
	row.NonComparableReason = "non_comparable_missing_window"
	summary.NonComparableCount++
	summary.NonComparableReason[row.NonComparableReason]++
}
