package telemetry

import "strings"

func escapeLabelValue(value string) string {
	// Escape according to Prometheus text label-value rules.
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "\n", "\\n")
	value = strings.ReplaceAll(value, "\r", "\\r")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	return value
}
