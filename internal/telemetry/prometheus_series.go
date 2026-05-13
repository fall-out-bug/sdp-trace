package telemetry

import "github.com/fall_out_bug/sdp-trace/internal/posture"

func gauge(name, help string, labels map[string]string, value float64) Series {
	copied := map[string]string{}
	for key, value := range labels {
		// Copy labels so later caller mutations cannot alter a built series.
		copied[key] = value
	}
	return Series{Name: name, Help: help, Type: "gauge", Labels: copied, Value: value}
}

func buildMetricSeries(rows []posture.MetricRow) ([]Series, error) {
	series := make([]Series, 0, len(rows)*3)
	for _, row := range rows {
		base, err := metricLabels(row)
		if err != nil {
			// Unsafe dimension labels block the whole telemetry document.
			return nil, err
		}
		series = appendMetricGauges(series, base, row)
	}
	return series, nil
}

func appendMetricGauges(series []Series, labels map[string]string, row posture.MetricRow) []Series {
	// Numerator, denominator, and not_assessed are separate gauges so dashboards
	// cannot infer hidden health scores.
	return append(series,
		gauge("sdp_trace_posture_metric_numerator", "Posture metric numerator from sdp-trace evidence posture export.", labels, float64(row.Numerator)),
		gauge("sdp_trace_posture_metric_denominator", "Posture metric denominator from sdp-trace evidence posture export.", labels, float64(row.Denominator)),
		gauge("sdp_trace_posture_metric_not_assessed", "Posture metric not assessed row count from sdp-trace evidence posture export.", labels, float64(row.NotAssessedCount)),
	)
}
