package posture

import (
	"github.com/fall_out_bug/sdp-trace/internal/query"
)

func metricMatches(metricID string, row query.QueryPackRow, signal PostureSignal, hasSignal bool) bool {
	// metricMatches keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	if expectedState, ok := rowStateMetrics[metricID]; ok {
		return row.EvidenceState == expectedState
	}
	return nonRowStateMetricMatches(metricID, row, signal, hasSignal)
}

func nonRowStateMetricMatches(metricID string, row query.QueryPackRow, signal PostureSignal, hasSignal bool) bool {
	// nonRowStateMetricMatches keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	if metricID == "unsupported_observer_rows" {
		return unsupportedObserverMetricMatches(row, signal, hasSignal)
	}
	return hasSignal && signalMetricMatches(metricID, signal)
}

func unsupportedObserverMetricMatches(row query.QueryPackRow, signal PostureSignal, hasSignal bool) bool {
	return row.EvidenceState == query.RowStateUnsupported || (hasSignal && signal.ObserverState == "unsupported")
}

func signalMetricMatches(metricID string, signal PostureSignal) bool {
	matches, ok := signalMetricPredicates[metricID]
	return ok && matches(signal)
}

// metricNotAssessed applies the evidence-boundary rule for not_assessed counts.
// For row_state metrics, follows row.EvidenceState. For signal metrics, signal absence is not_assessed.
func metricNotAssessed(def metricDef, row query.QueryPackRow, hasSignal bool) bool {
	// metricNotAssessed keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	if def.source == "posture_signal" {

		return !hasSignal
	}
	return row.EvidenceState == query.RowStateNotAssessed
}
