package harnessobs

import (
	"fmt"
)

func rejectUnsafeEvent(raw map[string]any, lineNo int) error {
	// rejectUnsafeEvent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if unsafeField, reason := findUnsafe(raw); unsafeField != "" {

		return fmt.Errorf("source line %d: unsafe_input:%s:%s", lineNo, unsafeField, reason)
	}
	return nil
}
