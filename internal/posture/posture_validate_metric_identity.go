package posture

func malformedMetricIdentity(row MetricRow) bool {
	return row.ID == "" || !validMetricID(row.MetricID) || row.MetricVersion != ProfileVer || !safeLabel(row.TimeWindow)
}
