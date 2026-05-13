package posture

import (
	"github.com/fall_out_bug/sdp-trace/internal/query"
)

func applyMetricCount(counts *metricCount, def metricDef, row query.QueryPackRow, signal PostureSignal, hasSignal bool) {
	// applyMetricCount keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	if metricMatches(def.id, row, signal, hasSignal) {
		counts.numerator++
	}
	if metricNotAssessed(def, row, hasSignal) {
		counts.notAssessed++
	}
	if missingPostureSignalMetric(def, hasSignal) {
		counts.sourceFieldState = "not_assessed"
	}
}

// missingPostureSignalMetric gates the evidence boundary for signal-sourced metrics.
// Absence of posture signal marks evidence-absent, downgrading sourceFieldState.
func missingPostureSignalMetric(def metricDef, hasSignal bool) bool {
	return def.source == "posture_signal" && !hasSignal
}
