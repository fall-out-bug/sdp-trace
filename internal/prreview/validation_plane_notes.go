package prreview

import "fmt"

func boolCount(value bool) int {
	if value {
		return 1
	}
	return 0
}

func appendPlaneValidationNotes(result PlaneResult, reasons, nextActions *[]string) {
	if result.Reason != "" {
		*reasons = append(*reasons, fmt.Sprintf("%s:%s", result.Plane, result.Status))
	}
	if result.NextAction != "" && !result.Usable {
		*nextActions = append(*nextActions, result.NextAction)
	}
}
