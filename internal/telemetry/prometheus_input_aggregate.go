package telemetry

import "github.com/fall_out_bug/sdp-trace/internal/posture"

func aggregateInputs(rows []posture.InputSelection) []Series {
	// Input rows report inventory distribution only; they do not score source
	// trust or expose selected file content.
	counts := map[string]Series{}
	for _, row := range rows {
		// Repository and window labels are retained as aggregate dimensions.
		countAggregate(counts, inputLabels(row), inputMetricName, inputMetricHelp)
	}
	return sortedAggregateValues(counts)
}

func inputLabels(row posture.InputSelection) map[string]string {
	// Input selection aggregation exposes trust-state distribution without
	// retaining raw source details.
	return map[string]string{
		"input_trust_state": row.InputTrustState,
		"repo":              row.Repository,
		"time_window":       row.TimeWindow,
	}
}
