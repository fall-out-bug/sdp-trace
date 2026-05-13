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

	refusalMetricName = "sdp_trace_posture_refusal"
	refusalMetricHelp = "Posture refusal fact count from sdp-trace evidence posture export."
	inputMetricName   = "sdp_trace_posture_input"
	inputMetricHelp   = "Posture selected input count from sdp-trace evidence posture export."
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

func renderPrometheusSeries(series []Series) string {
	var out strings.Builder
	currentFamily := ""
	for _, item := range series {
		if item.Name != currentFamily {
			// HELP/TYPE headers are emitted once per metric family.
			writePrometheusFamilyHeader(&out, item)
			currentFamily = item.Name
		}
		// Samples are rendered only after all series have passed label safety and
		// duplicate checks.
		writePrometheusSample(&out, item)
	}
	out.WriteString("# EOF\n")
	return out.String()
}

func writePrometheusFamilyHeader(out *strings.Builder, item Series) {
	// Help and type strings are fixed by the exporter, not sourced from labels.
	// Prometheus metadata therefore cannot carry posture input content.
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
}

func writePrometheusSample(out *strings.Builder, item Series) {
	// Values are rendered with Go's shortest decimal form to avoid synthetic
	// precision changes in diffs.
	out.WriteString(item.Name)
	out.WriteString(renderLabels(item.Labels))
	out.WriteByte(' ')
	out.WriteString(strconv.FormatFloat(item.Value, 'f', -1, 64))
	out.WriteByte('\n')
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

func finalizeSeries(series []Series) ([]Series, error) {
	for _, item := range series {
		if err := validateLabels(item.Labels); err != nil {
			// Reject unsafe labels before rendering any partial output.
			return nil, err
		}
	}
	return checkedSeries(series)
}

func checkedSeries(series []Series) ([]Series, error) {
	if err := rejectDuplicateSeries(series); err != nil {
		return nil, err
	}
	if len(series) > MaxSeries {
		// Bound exported cardinality so one posture file cannot create an
		// unbounded scrape surface.
		return nil, fmt.Errorf("series limit exceeded")
	}
	return series, nil
}

func aggregateRefusals(rows []posture.RefusalRow) []Series {
	// Refusal rows collapse to count gauges; raw refusal sources never leave the
	// posture export boundary.
	counts := map[string]Series{}
	for _, row := range rows {
		// Each row increments exactly one rendered label tuple.
		countAggregate(counts, refusalLabels(row), refusalMetricName, refusalMetricHelp)
	}
	return sortedAggregateValues(counts)
}

func refusalLabels(row posture.RefusalRow) map[string]string {
	// Refusal aggregation exposes state counts, not repository-specific source
	// paths or raw query-pack details.
	labels := map[string]string{
		"refusal_reason":    row.RefusalReason,
		"input_trust_state": row.InputTrustState,
	}
	if row.TimeWindow != "" {
		// Time window is optional for refusal rows and omitted when absent to keep
		// label cardinality stable.
		labels["time_window"] = row.TimeWindow
	}
	return labels
}

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

func countAggregate(counts map[string]Series, labels map[string]string, name, help string) {
	key := renderLabels(labels)
	item := counts[key]
	if item.Labels == nil {
		// First occurrence establishes the exact label set for the aggregate.
		item = gauge(name, help, labels, 0)
	}
	// Aggregates count rows; they do not weight by repository size or risk.
	item.Value++
	counts[key] = item
}

func sortedAggregateValues(counts map[string]Series) []Series {
	keys := make([]string, 0, len(counts))
	for key := range counts {
		// Rendered label tuples are the stable aggregate identity.
		keys = append(keys, key)
	}
	// Aggregates are sorted by rendered labels for deterministic output.
	sort.Strings(keys)
	out := make([]Series, 0, len(keys))
	for _, key := range keys {
		out = append(out, counts[key])
	}
	return out
}

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

func gauge(name, help string, labels map[string]string, value float64) Series {
	copied := map[string]string{}
	for key, value := range labels {
		// Copy labels so later caller mutations cannot alter a built series.
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
		if err := validateLabel(key, value); err != nil {
			// Return the first unsafe label to keep diagnostics deterministic.
			return err
		}
	}
	return nil
}

func validateLabel(key, value string) error {
	if !allowedLabelName(key) {
		// Closed label vocabulary keeps the telemetry contract small.
		return fmt.Errorf("unsupported label name: %s", key)
	}
	if value == "" {
		// Empty values are omitted at render time and are not unsafe.
		return nil
	}
	if unsafeLabelValue(value) {
		return fmt.Errorf("unsafe label value for key: %s", key)
	}
	return nil
}

func unsafeLabelValue(value string) bool {
	// Value safety combines size, encoding, and secret/path marker checks.
	return len(value) > MaxLabelValueBytes || !utf8.ValidString(value) || unsafeValue(value)
}

func allowedLabelName(value string) bool {
	switch value {
	case "metric_id", "metric_version", "dimension_key", "time_window", "repo", "team", "service", "harness", "change_type", "input_trust_state", "refusal_reason", "comparable", "non_comparable_reason":
		// Keep label names stable for Prometheus consumers.
		return true
	default:
		return false
	}
}

func unsafeValue(value string) bool {
	lower := strings.ToLower(value)
	// Lowercase secret-like markers and raw path/contact markers are checked
	// separately so case-insensitive tokens and raw separators are both caught.
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
			// First marker hit is enough to reject the label value.
			return true
		}
	}
	return false
}

func sortSeries(series []Series) {
	sort.Slice(series, func(i, j int) bool {
		left := series[i].Name + renderLabels(series[i].Labels)
		right := series[j].Name + renderLabels(series[j].Labels)
		// Name plus rendered labels is the Prometheus series identity.
		return left < right
	})
}

func rejectDuplicateSeries(series []Series) error {
	seen := map[string]struct{}{}
	for _, item := range series {
		key := item.Name + renderLabels(item.Labels)
		if _, ok := seen[key]; ok {
			// Duplicate series would be ambiguous to Prometheus scrapers.
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
	keys := nonEmptyLabelKeys(labels)
	if len(keys) == 0 {
		// All-empty label maps render as no label block.
		return ""
	}
	sort.Strings(keys)
	return renderSortedLabels(labels, keys)
}

func renderSortedLabels(labels map[string]string, keys []string) string {
	var out strings.Builder
	out.WriteByte('{')
	for i, key := range keys {
		if i > 0 {
			// Commas are written only between labels to keep text output
			// Prometheus-compatible.
			out.WriteByte(',')
		}
		// Keys are pre-sorted by renderLabels.
		// Values were validated before rendering and are escaped for text format.
		out.WriteString(key)
		out.WriteString("=\"")
		out.WriteString(escapeLabelValue(labels[key]))
		out.WriteByte('"')
	}
	out.WriteByte('}')
	return out.String()
}

func nonEmptyLabelKeys(labels map[string]string) []string {
	keys := make([]string, 0, len(labels))
	for key, value := range labels {
		if value != "" {
			// Empty label values are omitted rather than rendered as empty labels.
			keys = append(keys, key)
		}
	}
	return keys
}

func escapeLabelValue(value string) string {
	// Escape according to Prometheus text label-value rules.
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "\n", "\\n")
	value = strings.ReplaceAll(value, "\r", "\\r")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	return value
}
