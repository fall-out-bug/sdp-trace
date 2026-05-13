package telemetry

import "github.com/fall_out_bug/sdp-trace/internal/posture"

func comparableValue(value bool) float64 {
	if value {
		return 1
	}
	return 0
}

func buildMovementSeries(rows []posture.MovementRow) ([]Series, error) {
	series := make([]Series, 0, len(rows)*4)
	for _, row := range rows {
		base, err := movementLabels(row)
		if err != nil {
			// Movement labels share the same safety contract as metric labels.
			return nil, err
		}
		series = appendMovementGauges(series, base, row)
	}
	return series, nil
}

func appendMovementGauges(series []Series, labels map[string]string, row posture.MovementRow) []Series {
	// Movement comparability is explicit so consumers can avoid comparing
	// incomparable windows.
	return append(series,
		gauge("sdp_trace_posture_movement_current", "Current posture movement value from sdp-trace evidence posture export.", labels, float64(row.CurrentValue)),
		gauge("sdp_trace_posture_movement_previous", "Previous posture movement value from sdp-trace evidence posture export.", labels, float64(row.PreviousValue)),
		gauge("sdp_trace_posture_movement_delta", "Signed posture movement delta from sdp-trace evidence posture export.", labels, float64(row.Delta)),
		gauge("sdp_trace_posture_movement_comparable", "Posture movement comparability fact from sdp-trace evidence posture export.", labels, comparableValue(row.Comparable)),
	)
}
