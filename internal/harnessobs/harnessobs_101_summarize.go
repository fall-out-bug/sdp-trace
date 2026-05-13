package harnessobs

import (
	"fmt"

	"strings"
)

func Summarize(validation Validation) string {
	// Summarize keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var b strings.Builder

	fmt.Fprintf(&b, "Harness observation: %s (%s)\n", validation.ValidationState, validation.ReasonCode)
	fmt.Fprintf(&b, "Profile: %s\n", validation.ProfileID)
	fmt.Fprintf(&b, "Event schema: %s\n", validation.EventSchemaVersion)
	fmt.Fprintf(&b, "Events: %d\n", validation.EventCount)
	fmt.Fprintln(&b, "Dimensions:")
	for _, dim := range validation.Dimensions {
		writeSummaryDimension(&b, dim)
	}
	fmt.Fprintf(&b, "Boundary: %s\n", nonAuthority())
	return b.String()
}
