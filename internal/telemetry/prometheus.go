package telemetry

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/fall_out_bug/sdp-trace/internal/posture"
)

const (
	ProfilePrometheusTextV1 = "prometheus-text-v1"
	MaxSeries               = 10000
	MaxLabelValueBytes      = 1024
)

type Series struct {
	Name   string
	Help   string
	Type   string
	Labels map[string]string
	Value  float64
}

func RenderPrometheusText(result posture.ExportResult) (string, error) {
	return Render(result, ProfilePrometheusTextV1)
}

func Render(result posture.ExportResult, profile string) (string, error) {
	if profile != ProfilePrometheusTextV1 {
		return "", fmt.Errorf("unsupported telemetry profile")
	}
	return RenderPrometheus(result)
}

func RenderPrometheus(result posture.ExportResult) (string, error) {
	if err := posture.ValidateExportResult(result); err != nil {
		return "", err
	}
	series, err := BuildSeries(result)
	if err != nil {
		return "", err
	}
	if len(series) == 0 {
		return "# sdp_trace_posture no_rows\n# EOF\n", nil
	}
	sortSeries(series)
	var out strings.Builder
	currentFamily := ""
	for _, item := range series {
		if item.Name != currentFamily {
			out.WriteString("# HELP ")
			out.WriteString(item.Name)
			out.WriteByte(' ')
			out.WriteString(item.Help)
			out.WriteByte('\n')
			out.WriteString("# TYPE ")
			out.WriteString(item.Name)
			out.WriteByte(' ')
			out.WriteString(item.Type)
			out.WriteByte('\n')
			currentFamily = item.Name
		}
		out.WriteString(item.Name)
		out.WriteString(renderLabels(item.Labels))
		out.WriteByte(' ')
		out.WriteString(strconv.FormatFloat(item.Value, 'f', -1, 64))
		out.WriteByte('\n')
	}
	out.WriteString("# EOF\n")
	return out.String(), nil
}

func BuildSeries(result posture.ExportResult) ([]Series, error) {
	var series []Series
	for _, row := range result.MetricRows {
		base, err := metricLabels(row)
		if err != nil {
			return nil, err
		}
		series = append(series,
			gauge("sdp_trace_posture_metric_numerator", "Posture metric numerator from sdp-trace evidence posture export.", base, float64(row.Numerator)),
			gauge("sdp_trace_posture_metric_denominator", "Posture metric denominator from sdp-trace evidence posture export.", base, float64(row.Denominator)),
			gauge("sdp_trace_posture_metric_not_assessed", "Posture metric not assessed row count from sdp-trace evidence posture export.", base, float64(row.NotAssessedCount)),
		)
	}
	for _, row := range result.MovementRows {
		base, err := movementLabels(row)
		if err != nil {
			return nil, err
		}
		series = append(series,
			gauge("sdp_trace_posture_movement_current", "Current posture movement value from sdp-trace evidence posture export.", base, float64(row.CurrentValue)),
			gauge("sdp_trace_posture_movement_previous", "Previous posture movement value from sdp-trace evidence posture export.", base, float64(row.PreviousValue)),
			gauge("sdp_trace_posture_movement_delta", "Signed posture movement delta from sdp-trace evidence posture export.", base, float64(row.Delta)),
			gauge("sdp_trace_posture_movement_comparable", "Posture movement comparability fact from sdp-trace evidence posture export.", base, comparableValue(row.Comparable)),
		)
	}
	series = append(series, aggregateRefusals(result.RefusalRows)...)
	series = append(series, aggregateInputs(result.InputSelection)...)
	for _, item := range series {
		if err := validateLabels(item.Labels); err != nil {
			return nil, err
		}
	}
	if err := rejectDuplicateSeries(series); err != nil {
		return nil, err
	}
	if len(series) > MaxSeries {
		return nil, fmt.Errorf("series limit exceeded")
	}
	return series, nil
}

func aggregateRefusals(rows []posture.RefusalRow) []Series {
	counts := map[string]Series{}
	for _, row := range rows {
		labels := map[string]string{
			"refusal_reason":    row.RefusalReason,
			"input_trust_state": row.InputTrustState,
		}
		if row.TimeWindow != "" {
			labels["time_window"] = row.TimeWindow
		}
		key := renderLabels(labels)
		item := counts[key]
		if item.Labels == nil {
			item = gauge("sdp_trace_posture_refusal", "Posture refusal fact count from sdp-trace evidence posture export.", labels, 0)
		}
		item.Value++
		counts[key] = item
	}
	return sortedAggregateValues(counts)
}

func aggregateInputs(rows []posture.InputSelection) []Series {
	counts := map[string]Series{}
	for _, row := range rows {
		labels := map[string]string{
			"input_trust_state": row.InputTrustState,
			"repo":              row.Repository,
			"time_window":       row.TimeWindow,
		}
		key := renderLabels(labels)
		item := counts[key]
		if item.Labels == nil {
			item = gauge("sdp_trace_posture_input", "Posture selected input count from sdp-trace evidence posture export.", labels, 0)
		}
		item.Value++
		counts[key] = item
	}
	return sortedAggregateValues(counts)
}

func sortedAggregateValues(counts map[string]Series) []Series {
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := make([]Series, 0, len(keys))
	for _, key := range keys {
		out = append(out, counts[key])
	}
	return out
}

func metricLabels(row posture.MetricRow) (map[string]string, error) {
	labels := map[string]string{
		"metric_id":      row.MetricID,
		"metric_version": row.MetricVersion,
		"dimension_key":  row.DimensionKey,
		"time_window":    row.TimeWindow,
	}
	for _, key := range []string{"repo", "team", "service", "harness", "change_type"} {
		if value := row.Dimensions[key]; value != "" {
			labels[key] = value
		}
	}
	return labels, validateLabels(labels)
}

func movementLabels(row posture.MovementRow) (map[string]string, error) {
	labels := map[string]string{
		"metric_id":      row.MetricID,
		"metric_version": row.MetricVersion,
		"dimension_key":  row.DimensionKey,
		"comparable":     strconv.FormatBool(row.Comparable),
	}
	if row.NonComparableReason != "" {
		labels["non_comparable_reason"] = row.NonComparableReason
	}
	return labels, validateLabels(labels)
}

func gauge(name, help string, labels map[string]string, value float64) Series {
	copied := map[string]string{}
	for key, value := range labels {
		copied[key] = value
	}
	return Series{Name: name, Help: help, Type: "gauge", Labels: copied, Value: value}
}

func comparableValue(value bool) float64 {
	if value {
		return 1
	}
	return 0
}

func validateLabels(labels map[string]string) error {
	for key, value := range labels {
		if !allowedLabelName(key) {
			return fmt.Errorf("unsupported label name: %s", key)
		}
		if value == "" {
			continue
		}
		if len(value) > MaxLabelValueBytes || !utf8.ValidString(value) || unsafeValue(value) {
			return fmt.Errorf("unsafe label value for key: %s", key)
		}
	}
	return nil
}

func allowedLabelName(value string) bool {
	switch value {
	case "metric_id", "metric_version", "dimension_key", "time_window", "repo", "team", "service", "harness", "change_type", "input_trust_state", "refusal_reason", "comparable", "non_comparable_reason":
		return true
	default:
		return false
	}
}

func unsafeValue(value string) bool {
	lower := strings.ToLower(value)
	return containsAnyMarker(lower, unsafeLowerMarkers) || containsAnyMarker(value, unsafeRawMarkers)
}

var unsafeLowerMarkers = []string{
	"http://",
	"https://",
	"secret",
	"token",
	"credential",
	"password",
	"bearer",
	"api_key",
	"access_key",
	"private",
}

var unsafeRawMarkers = []string{"@", "/", "\\"}

func containsAnyMarker(value string, markers []string) bool {
	for _, marker := range markers {
		if strings.Contains(value, marker) {
			return true
		}
	}
	return false
}

func sortSeries(series []Series) {
	sort.Slice(series, func(i, j int) bool {
		left := series[i].Name + renderLabels(series[i].Labels)
		right := series[j].Name + renderLabels(series[j].Labels)
		return left < right
	})
}

func rejectDuplicateSeries(series []Series) error {
	seen := map[string]struct{}{}
	for _, item := range series {
		key := item.Name + renderLabels(item.Labels)
		if _, ok := seen[key]; ok {
			return fmt.Errorf("duplicate series for metric %s", item.Name)
		}
		seen[key] = struct{}{}
	}
	return nil
}

func renderLabels(labels map[string]string) string {
	if len(labels) == 0 {
		return ""
	}
	keys := make([]string, 0, len(labels))
	for key, value := range labels {
		if value != "" {
			keys = append(keys, key)
		}
	}
	if len(keys) == 0 {
		return ""
	}
	sort.Strings(keys)
	var out strings.Builder
	out.WriteByte('{')
	for i, key := range keys {
		if i > 0 {
			out.WriteByte(',')
		}
		out.WriteString(key)
		out.WriteString("=\"")
		out.WriteString(escapeLabelValue(labels[key]))
		out.WriteByte('"')
	}
	out.WriteByte('}')
	return out.String()
}

func escapeLabelValue(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "\n", "\\n")
	value = strings.ReplaceAll(value, "\r", "\\r")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	return value
}
