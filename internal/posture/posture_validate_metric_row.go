package posture

import (
	"fmt"
)

func validateMetricRow(row MetricRow) error {
	// validateMetricRow keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	if err := validateMetricRowShape(row); err != nil {
		return err
	}
	if err := validateMetricDimensions(row.Dimensions); err != nil {
		return err
	}
	return validateMetricTrustSummary(row.InputTrustStateSummary)
}

func validateMetricRowShape(row MetricRow) error {
	// validateMetricRowShape keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	if malformedMetricIdentity(row) || malformedMetricCounts(row) || malformedMetricSource(row) {
		return fmt.Errorf("malformed posture export metric_row")
	}
	return nil
}

func malformedMetricIdentity(row MetricRow) bool {
	return row.ID == "" || !validMetricID(row.MetricID) || row.MetricVersion != ProfileVer || !safeLabel(row.TimeWindow)
}

func malformedMetricCounts(row MetricRow) bool {
	return row.Numerator < 0 || row.Denominator < 0 || row.Unit != "rows" || row.NotAssessedCount < 0
}
