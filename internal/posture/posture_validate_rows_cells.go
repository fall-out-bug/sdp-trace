package posture

import (
	"fmt"
)

func validateExportCollections(result ExportResult) error {
	// validateExportCollections keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	if !validExportGrouping(result.GroupingSetID, result.ActiveGroupingKeys) {
		return fmt.Errorf("malformed posture export grouping")
	}

	if !hasRequiredCollections(result) {
		return fmt.Errorf("malformed posture export missing required collection")
	}
	if !hasOutputSafety(result.OutputSafety.VerifiedAbsentSensitiveClasses) {
		return fmt.Errorf("malformed posture export output_safety")
	}
	return nil
}

func validExportGrouping(groupingSet string, keys []string) bool {
	return len(groupingKeys(groupingSet)) > 0 && len(keys) >= 2
}

func hasRequiredCollections(result ExportResult) bool {
	// hasRequiredCollections keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	for _, present := range []bool{
		result.InputSelection != nil,
		result.MetricRows != nil,
		result.MovementRows != nil,
		result.RefusalRows != nil,
		result.Handoff != nil,
		result.MovementSummary.NonComparableReason != nil,
	} {
		if !present {

			return false
		}
	}
	return true
}

func hasOutputSafety(classes []string) bool {
	return len(classes) > 0
}
