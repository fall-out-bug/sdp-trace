package telemetry

import "github.com/fall_out_bug/sdp-trace/internal/posture"

type Series struct {
	Name   string
	Help   string
	Type   string
	Labels map[string]string
	Value  float64
}

func BuildSeries(result posture.ExportResult) ([]Series, error) {
	metricSeries, err := buildMetricSeries(result.MetricRows)
	if err != nil {
		return nil, err
	}
	// Movement and metric families are built separately so dashboards cannot
	// treat movement deltas as gate verdicts.
	movementSeries, err := buildMovementSeries(result.MovementRows)
	if err != nil {
		return nil, err
	}
	// Build all families before final validation so duplicate series across
	// sources are rejected together.
	series := make([]Series, 0, len(metricSeries)+len(movementSeries)+len(result.RefusalRows)+len(result.InputSelection))
	// Series are appended by family so final sorting can give deterministic
	// output without changing the family-specific construction rules.
	series = append(series, metricSeries...)
	series = append(series, movementSeries...)
	series = append(series, aggregateRefusals(result.RefusalRows)...)
	series = append(series, aggregateInputs(result.InputSelection)...)
	return finalizeSeries(series)
}
