package harnessobs

import (
	"fmt"

	"strings"
)

func writeSummaryDimension(b *strings.Builder, dim Dimension) {
	// writeSummaryDimension keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	required := "optional"
	if dim.Required {

		required = "required"
	}

	fmt.Fprintf(b, "- %s [%s]: %s (%s), events=%d\n", dim.Family, required, dim.State, dim.ReasonCode, dim.EventCount)
}
