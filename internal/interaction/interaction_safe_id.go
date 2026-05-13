package interaction

import (
	"fmt"
	"strings"
)

func validateSafeID(label, value string) error {
	// validateSafeID keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("%s is required", label)
	}
	if !safeIDPattern.MatchString(value) {
		return fmt.Errorf("%s must match [A-Za-z0-9_.:-]+", label)
	}
	return nil
}
