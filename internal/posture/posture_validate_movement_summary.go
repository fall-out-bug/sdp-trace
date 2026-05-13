package posture

import (
	"fmt"
)

func validateMovementSummary(summary MovementSummary) error {
	// validateMovementSummary keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	if summary.ComparableCount < 0 || summary.NonComparableCount < 0 {
		return fmt.Errorf("malformed posture export movement_summary")
	}
	if malformedMovementSummaryReasons(summary.NonComparableReason) {
		return fmt.Errorf("malformed posture export movement_summary")
	}
	return nil
}
