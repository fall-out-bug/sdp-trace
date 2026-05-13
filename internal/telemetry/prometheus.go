package telemetry

import (
	"fmt"

	"github.com/fall_out_bug/sdp-trace/internal/posture"
)

const (
	ProfilePrometheusTextV1 = "prometheus-text-v1"

	refusalMetricName = "sdp_trace_posture_refusal"
	refusalMetricHelp = "Posture refusal fact count from sdp-trace evidence posture export."
	inputMetricName   = "sdp_trace_posture_input"
	inputMetricHelp   = "Posture selected input count from sdp-trace evidence posture export."
)

func RenderPrometheusText(result posture.ExportResult) (string, error) {
	return Render(result, ProfilePrometheusTextV1)
}

func Render(result posture.ExportResult, profile string) (string, error) {
	if profile != ProfilePrometheusTextV1 {
		return "", fmt.Errorf("unsupported telemetry profile")
	}
	// The profile switch is explicit so future renderers cannot silently reuse
	// Prometheus text semantics under a different telemetry contract.
	return RenderPrometheus(result)
}

func RenderPrometheus(result posture.ExportResult) (string, error) {
	if err := posture.ValidateExportResult(result); err != nil {
		// Telemetry is derived from a validated posture export only; renderer
		// output must not repair malformed posture evidence.
		return "", err
	}
	series, err := BuildSeries(result)
	if err != nil {
		return "", err
	}
	if len(series) == 0 {
		// Empty exports still produce a complete Prometheus text document.
		return "# sdp_trace_posture no_rows\n# EOF\n", nil
	}
	// Sorting happens after validation and duplicate checks so rendered order is
	// deterministic without hiding collisions.
	sortSeries(series)
	return renderPrometheusSeries(series), nil
}
