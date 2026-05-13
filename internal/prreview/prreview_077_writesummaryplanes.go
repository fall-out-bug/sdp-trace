package prreview

import (
	"fmt"

	"strings"
)

func writeSummaryPlanes(b *strings.Builder, planes []PlaneResult) {
	// writeSummaryPlanes keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if len(planes) > 0 {
		b.WriteString("\nPlanes\n")
		for _, plane := range planes {

			fmt.Fprintf(b, "- %s: %s", plane.Plane, plane.Status)
			if plane.NextAction != "" {
				fmt.Fprintf(b, " next_action=%s", safeText(plane.NextAction))
			}
			b.WriteString("\n")
		}
	}
}
