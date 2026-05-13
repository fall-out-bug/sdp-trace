package posture

func malformedMovementIdentity(row MovementRow) bool {
	return row.ID == "" || !validMetricID(row.MetricID) || row.MetricVersion != ProfileVer || row.DimensionKey == ""
}
