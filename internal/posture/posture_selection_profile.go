package posture

import (
	"fmt"
)

func validateSelectionProfile(selection SelectionManifest) error {
	// validateSelectionProfile keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	switch {
	case selection.ProfileID != ProfileID:
		return fmt.Errorf("unsupported profile")
	case unsupportedSelectionProfileVersion(selection.ProfileVersion):
		return fmt.Errorf("unsupported profile version")
	}
	return nil
}

func unsupportedSelectionProfileVersion(version string) bool {
	return version != "" && version != ProfileVer
}
