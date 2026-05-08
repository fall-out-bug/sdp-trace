package telemetry

import (
	"strings"
	"testing"

	"github.com/fall_out_bug/sdp-trace/internal/posture"
)

func TestRenderPrometheusTextRendersFamiliesAndEOF(t *testing.T) {
	result := validResult()
	out, err := RenderPrometheusText(result)
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	for _, want := range []string{
		"# HELP sdp_trace_posture_input ",
		"# TYPE sdp_trace_posture_input gauge\n",
		"# HELP sdp_trace_posture_metric_numerator ",
		"# TYPE sdp_trace_posture_metric_numerator gauge\n",
		`sdp_trace_posture_metric_numerator{dimension_key="repo=repo-a",metric_id="cannot_verify_rows",metric_version="v1",repo="repo-a",time_window="2026-w02"} 2`,
		`sdp_trace_posture_metric_not_assessed{dimension_key="repo=repo-a",metric_id="cannot_verify_rows",metric_version="v1",repo="repo-a",time_window="2026-w02"} 3`,
		`sdp_trace_posture_movement_comparable{comparable="false",dimension_key="repo=repo-a",metric_id="cannot_verify_rows",metric_version="v1",non_comparable_reason="non_comparable_missing_window"} 0`,
		`sdp_trace_posture_refusal{input_trust_state="cannot_verify_input",refusal_reason="malformed_input",time_window="2026-w02"} 1`,
		"# EOF\n",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in:\n%s", want, out)
		}
	}
}

func TestRenderPrometheusTextOrdersByMetricThenLabels(t *testing.T) {
	result := validResult()
	result.MetricRows = append(result.MetricRows, posture.MetricRow{
		MetricID:      "cannot_verify_rows",
		MetricVersion: "v1",
		Numerator:     1,
		Denominator:   3,
		TimeWindow:    "2026-w02",
		Dimensions:    map[string]string{"repo": "repo-0"},
		DimensionKey:  "repo=repo-0",
	})
	out, err := RenderPrometheusText(result)
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	first := strings.Index(out, `sdp_trace_posture_metric_numerator{dimension_key="repo=repo-0"`)
	second := strings.Index(out, `sdp_trace_posture_metric_numerator{dimension_key="repo=repo-a"`)
	if first == -1 || second == -1 || first > second {
		t.Fatalf("unexpected ordering:\n%s", out)
	}
}

func TestRenderPrometheusTextEscapesLabelValues(t *testing.T) {
	result := validResult()
	result.MetricRows[0].Dimensions["repo"] = "repo\"a\nb\rc"
	result.MetricRows[0].DimensionKey = "repo=\"a\nb\rc"
	out, err := RenderPrometheusText(result)
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if !strings.Contains(out, `dimension_key="repo=\"a\nb\rc"`) || !strings.Contains(out, `repo="repo\"a\nb\rc"`) {
		t.Fatalf("labels were not escaped:\n%s", out)
	}
}

func TestRenderPrometheusTextRejectsDuplicateSeries(t *testing.T) {
	result := validResult()
	result.MetricRows = append(result.MetricRows, result.MetricRows[0])
	if out, err := RenderPrometheusText(result); err == nil || out != "" {
		t.Fatalf("render = %q, %v; want empty duplicate-series error", out, err)
	}
}

func TestRenderPrometheusTextRejectsUnsafeLabelsWithoutPartialOutput(t *testing.T) {
	result := validResult()
	result.MetricRows[0].Dimensions["repo"] = "https://example.test/private"
	if out, err := RenderPrometheusText(result); err == nil || out != "" {
		t.Fatalf("render = %q, %v; want empty unsafe error", out, err)
	}
}

func TestRenderPrometheusTextRejectsLongLabelValue(t *testing.T) {
	result := validResult()
	result.MetricRows[0].Dimensions["repo"] = strings.Repeat("a", 1025)
	if out, err := RenderPrometheusText(result); err == nil || out != "" {
		t.Fatalf("render = %q, %v; want empty long-label error", out, err)
	}
}

func TestRenderPrometheusTextRejectsUnsupportedInput(t *testing.T) {
	result := validResult()
	result.SchemaVersion = "future"
	if out, err := RenderPrometheusText(result); err == nil || out != "" {
		t.Fatalf("render = %q, %v; want empty unsupported error", out, err)
	}
}

func TestRenderPrometheusTextEnforcesSeriesLimit(t *testing.T) {
	result := validResult()
	result.MovementRows = nil
	result.RefusalRows = nil
	result.InputSelection = nil
	result.MetricRows = make([]posture.MetricRow, 3334)
	for i := range result.MetricRows {
		result.MetricRows[i] = posture.MetricRow{
			MetricID:      "cannot_verify_rows",
			MetricVersion: "v1",
			TimeWindow:    "2026-w02",
			DimensionKey:  "repo=repo-a",
		}
	}
	if out, err := RenderPrometheusText(result); err == nil || out != "" {
		t.Fatalf("render = %q, %v; want empty series-limit error", out, err)
	}
}

func TestRenderPrometheusTextEmptyResult(t *testing.T) {
	result := validResult()
	result.MetricRows = nil
	result.MovementRows = nil
	result.RefusalRows = nil
	result.InputSelection = nil
	out, err := RenderPrometheusText(result)
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	if out != "# sdp_trace_posture no_rows\n# EOF\n" {
		t.Fatalf("empty output = %q", out)
	}
}

func validResult() posture.ExportResult {
	return posture.ExportResult{
		SchemaVersion:        posture.SchemaVersion,
		ExportProfileID:      posture.ProfileID,
		ExportProfileVersion: posture.ProfileVer,
		InputSelection: []posture.InputSelection{
			{Repository: "repo-a", TimeWindow: "2026-w02", InputTrustState: "trusted_input"},
			{Repository: "repo-a", TimeWindow: "2026-w02", InputTrustState: "trusted_input"},
		},
		MetricRows: []posture.MetricRow{
			{
				MetricID:         "cannot_verify_rows",
				MetricVersion:    "v1",
				Numerator:        2,
				Denominator:      5,
				NotAssessedCount: 3,
				TimeWindow:       "2026-w02",
				Dimensions:       map[string]string{"repo": "repo-a"},
				DimensionKey:     "repo=repo-a",
			},
		},
		MovementRows: []posture.MovementRow{
			{
				MetricID:            "cannot_verify_rows",
				MetricVersion:       "v1",
				DimensionKey:        "repo=repo-a",
				CurrentValue:        2,
				PreviousValue:       1,
				Delta:               1,
				Comparable:          false,
				NonComparableReason: "non_comparable_missing_window",
			},
		},
		RefusalRows: []posture.RefusalRow{
			{TimeWindow: "2026-w02", RefusalReason: "malformed_input", InputTrustState: "cannot_verify_input"},
		},
	}
}
