package posture

import (
	"fmt"
)

func readSelection(path string) (SelectionManifest, error) {
	return readJSONFile[SelectionManifest](path)
}

func validateSelection(selection SelectionManifest) error {
	// validateSelection keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	if selection.SchemaVersion != SelectionSchemaVersion {
		return fmt.Errorf("unsupported selection schema")
	}
	if err := validateSelectionProfile(selection); err != nil {
		return err
	}
	if err := validateSelectionGrouping(selection); err != nil {
		return err
	}
	return validateSelectionRepositories(selection)
}
