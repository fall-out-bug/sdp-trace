package posture

func malformedMetricCounts(row MetricRow) bool {
	return row.Numerator < 0 || row.Denominator < 0 || row.Unit != "rows" || row.NotAssessedCount < 0
}
