package posture

func malformedMetricSource(row MetricRow) bool {
	return missingMetricSourceRefs(row) || missingMetricTrustSource(row)
}

func missingMetricSourceRefs(row MetricRow) bool {
	return row.Dimensions == nil || row.DimensionKey == "" || row.SourceInputRefs == nil || row.SourceArtifactDigestSet == ""
}

func missingMetricTrustSource(row MetricRow) bool {
	return !validSourceFieldState(row.SourceFieldState) || row.InputTrustStateSummary == nil
}
