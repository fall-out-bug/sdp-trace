package telemetry

import (
	"strconv"
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
	row := result.MetricRows[0]
	row.ID = "metric.0002"
	row.Numerator = 1
	row.Denominator = 3
	row.Dimensions = map[string]string{"repo": "repo-0"}
	row.DimensionKey = "repo=repo-0"
	result.MetricRows = append(result.MetricRows, row)
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
	out := renderLabels(map[string]string{
		"dimension_key": "repo=\"a\nb\rc",
		"repo":          "repo\"a\nb\rc",
	})
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

func TestRenderPrometheusTextDoesNotLetDimensionsOverrideDimensionKey(t *testing.T) {
	result := validResult()
	result.MetricRows[0].Dimensions["dimension_key"] = "repo=spoofed"
	out, err := RenderPrometheusText(result)
	if err == nil || out != "" {
		t.Fatalf("render = %q, %v; want empty malformed-dimension error", out, err)
	}
}

func TestRenderPrometheusTextRejectsUnsafeLabelsWithoutPartialOutput(t *testing.T) {
	for _, value := range []string{
		"http://example.test",
		"https://example.test/private",
		"secret",
		"token",
		"credential",
		"password",
		"api_key",
		"access_key",
		"bearer",
		"private",
		"user@example.test",
		"owner/repo",
		`owner\repo`,
	} {
		result := validResult()
		result.MetricRows[0].Dimensions["repo"] = value
		if out, err := RenderPrometheusText(result); err == nil || out != "" {
			t.Fatalf("render = %q, %v; want empty unsafe error for %q", out, err, value)
		}
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
	result.MovementRows = []posture.MovementRow{}
	result.RefusalRows = []posture.RefusalRow{}
	result.InputSelection = []posture.InputSelection{}
	result.MetricRows = make([]posture.MetricRow, 3334)
	for i := range result.MetricRows {
		repo := "repo-" + strconv.Itoa(i+1)
		result.MetricRows[i] = posture.MetricRow{
			ID:                      "metric." + strings.Repeat("0", 4-len(strconv.Itoa(i+1))) + strconv.Itoa(i+1),
			MetricID:                "cannot_verify_rows",
			MetricVersion:           "v1",
			Numerator:               1,
			Denominator:             1,
			Unit:                    "rows",
			TimeWindow:              "2026-w02",
			Dimensions:              map[string]string{"repo": repo},
			DimensionKey:            "repo=" + repo,
			SourceInputRefs:         []string{},
			SourceArtifactDigestSet: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			SourceFieldState:        "present",
			InputTrustStateSummary:  map[string]int{},
		}
	}
	if out, err := RenderPrometheusText(result); err == nil || out != "" {
		t.Fatalf("render = %q, %v; want empty series-limit error", out, err)
	}
}

func TestRenderPrometheusTextEmptyResult(t *testing.T) {
	result := validResult()
	result.MetricRows = []posture.MetricRow{}
	result.MovementRows = []posture.MovementRow{}
	result.RefusalRows = []posture.RefusalRow{}
	result.InputSelection = []posture.InputSelection{}
	result.MovementSummary = posture.MovementSummary{NonComparableReason: map[string]int{}}
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
		ExportID:             "export:0123456789abcdef",
		Producer:             "sdp-trace",
		GeneratedAt:          "2026-01-10T00:00:00Z",
		GroupingSetID:        posture.GroupingRepoWindow,
		ActiveGroupingKeys:   []string{"repo", "time_window"},
		InputSelection: []posture.InputSelection{
			{InputID: "input-a", Repository: "repo-a", TimeWindow: "2026-w02", PathRedactedID: "artifact:query:not_assessed0000", InputTrustState: "trusted_input"},
			{InputID: "input-b", Repository: "repo-a", TimeWindow: "2026-w02", PathRedactedID: "artifact:query:not_assessed0000", InputTrustState: "trusted_input"},
		},
		MetricRows: []posture.MetricRow{
			{
				MetricID:                "cannot_verify_rows",
				MetricVersion:           "v1",
				ID:                      "metric.0001",
				Numerator:               2,
				Denominator:             5,
				Unit:                    "rows",
				NotAssessedCount:        3,
				TimeWindow:              "2026-w02",
				Dimensions:              map[string]string{"repo": "repo-a"},
				DimensionKey:            "repo=repo-a",
				SourceInputRefs:         []string{"input-a"},
				SourceArtifactDigestSet: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
				SourceFieldState:        "present",
				InputTrustStateSummary:  map[string]int{"trusted_input": 1},
			},
		},
		MovementRows: []posture.MovementRow{
			{
				ID:                   "movement.0001",
				MetricID:             "cannot_verify_rows",
				MetricVersion:        "v1",
				DimensionKey:         "repo=repo-a",
				CurrentMetricRowRef:  "metric.0001",
				PreviousMetricRowRef: "metric.0001",
				CurrentValue:         2,
				PreviousValue:        1,
				Delta:                1,
				ComparisonBasis:      "non_comparable_missing_window",
				Comparable:           false,
				NonComparableReason:  "non_comparable_missing_window",
			},
		},
		MovementSummary: posture.MovementSummary{
			ComparableCount:     0,
			NonComparableCount:  1,
			NonComparableReason: map[string]int{"non_comparable_missing_window": 1},
		},
		RefusalRows: []posture.RefusalRow{
			{ID: "refusal.0001", InputID: "input-a", TimeWindow: "2026-w02", RefusalReason: "malformed_input", InputTrustState: "cannot_verify_input"},
		},
		Handoff: map[string]string{},
		OutputSafety: posture.OutputSafety{
			VerifiedAbsentSensitiveClasses: posture.SensitiveClasses(),
		},
	}
}
