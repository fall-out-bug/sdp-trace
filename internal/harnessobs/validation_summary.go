package harnessobs

import (
	"fmt"
	"strings"
)

// Validation summary renders the portable human-readable observation verdict.
// The boundary line keeps the report explicit that it is not external proof.
func Summarize(validation Validation) string {
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

func writeSummaryDimension(b *strings.Builder, dim Dimension) {
	required := "optional"
	if dim.Required {
		required = "required"
	}
	fmt.Fprintf(b, "- %s [%s]: %s (%s), events=%d\n", dim.Family, required, dim.State, dim.ReasonCode, dim.EventCount)
}
