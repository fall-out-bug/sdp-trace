package posture

import (
	"fmt"
)

func validateMetricDimensions(dimensions map[string]string) error {
	// validateMetricDimensions keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	for key, value := range dimensions {

		if !validDimensionName(key) || !safeLabel(value) {
			return fmt.Errorf("malformed posture export metric_row dimensions")
		}
	}
	return nil
}

func validateMetricTrustSummary(summary map[string]int) error {
	// validateMetricTrustSummary keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	for state, count := range summary {

		if !validInputTrustState(state) || count < 0 {
			return fmt.Errorf("malformed posture export input_trust_state_summary")
		}
	}
	return nil
}
