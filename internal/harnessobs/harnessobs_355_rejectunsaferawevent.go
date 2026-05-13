package harnessobs

import (
	"fmt"
)

func rejectUnsafeRawEvent(raw map[string]any, lineNo int) error {
	// rejectUnsafeRawEvent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if unsafeField, reason := findUnsafeRawEvent(raw); unsafeField != "" {

		return fmt.Errorf("raw source line %d: unsafe_input:%s:%s", lineNo, unsafeField, reason)
	}
	return nil
}
