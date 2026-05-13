package telemetry

import (
	"strconv"

	"github.com/fall_out_bug/sdp-trace/internal/posture"
)

func metricLabels(row posture.MetricRow) (map[string]string, error) {
	// Metric labels carry the closed posture dimension vocabulary into
	// Prometheus without exposing arbitrary selection metadata.
	labels := map[string]string{
		"metric_id":      row.MetricID,
		"metric_version": row.MetricVersion,
		"dimension_key":  row.DimensionKey,
		"time_window":    row.TimeWindow,
	}
	for _, key := range []string{"repo", "team", "service", "harness", "change_type"} {
		if value := row.Dimensions[key]; value != "" {
			// Only the supported public dimensions become Prometheus labels.
			labels[key] = value
		}
	}
	return labels, validateLabels(labels)
}

func movementLabels(row posture.MovementRow) (map[string]string, error) {
	// Movement labels include comparability so consumers can separate missing
	// windows from real changes.
	labels := map[string]string{
		"metric_id":      row.MetricID,
		"metric_version": row.MetricVersion,
		"dimension_key":  row.DimensionKey,
		"comparable":     strconv.FormatBool(row.Comparable),
	}
	if row.NonComparableReason != "" {
		// Non-comparable reason is present only when the movement row supplies it.
		labels["non_comparable_reason"] = row.NonComparableReason
	}
	return labels, validateLabels(labels)
}
