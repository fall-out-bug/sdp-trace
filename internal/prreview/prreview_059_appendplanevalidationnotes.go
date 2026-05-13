package prreview

import (
	"fmt"
)

func appendPlaneValidationNotes(result PlaneResult, reasons, nextActions *[]string) {
	// appendPlaneValidationNotes keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if result.Reason != "" {

		*reasons = append(*reasons, fmt.Sprintf("%s:%s", result.Plane, result.Status))
	}
	if result.NextAction != "" && !result.Usable {

		*nextActions = append(*nextActions, result.NextAction)
	}
}
