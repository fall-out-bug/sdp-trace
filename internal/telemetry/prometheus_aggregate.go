package telemetry

import (
	"sort"

	"github.com/fall_out_bug/sdp-trace/internal/posture"
)

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
