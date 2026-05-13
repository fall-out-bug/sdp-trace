package posture

import (
	"fmt"
)

func validateSelectionGrouping(selection SelectionManifest) error {
	// validateSelectionGrouping keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	if len(groupingKeys(selection.GroupingSetID)) == 0 {
		return fmt.Errorf("unsupported grouping set")
	}

	if !groupingAllowedByExposure(selection.GroupingSetID, selection.DimensionExposurePolicy) {
		return fmt.Errorf("dimension exposure policy excludes grouping key")
	}
	return nil
}

func validateSelectionRepositories(selection SelectionManifest) error {
	// validateSelectionRepositories keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	if len(selection.Repositories) == 0 {
		return fmt.Errorf("empty selection")
	}
	return nil
}
