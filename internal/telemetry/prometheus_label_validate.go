package telemetry

import (
	"fmt"
	"unicode/utf8"
)

func validateLabels(labels map[string]string) error {
	for key, value := range labels {
		if err := validateLabel(key, value); err != nil {
			// Return the first unsafe label to keep diagnostics deterministic.
			return err
		}
	}
	return nil
}

func validateLabel(key, value string) error {
	if !allowedLabelName(key) {
		// Closed label vocabulary keeps the telemetry contract small.
		return fmt.Errorf("unsupported label name: %s", key)
	}
	if value == "" {
		// Empty values are omitted at render time and are not unsafe.
		return nil
	}
	if unsafeLabelValue(value) {
		return fmt.Errorf("unsafe label value for key: %s", key)
	}
	return nil
}

func unsafeLabelValue(value string) bool {
	// Value safety combines size, encoding, and secret/path marker checks.
	return len(value) > MaxLabelValueBytes || !utf8.ValidString(value) || unsafeValue(value)
}

func allowedLabelName(value string) bool {
	switch value {
	case "metric_id", "metric_version", "dimension_key", "time_window", "repo", "team", "service", "harness", "change_type", "input_trust_state", "refusal_reason", "comparable", "non_comparable_reason":
		// Keep label names stable for Prometheus consumers.
		return true
	default:
		return false
	}
}
