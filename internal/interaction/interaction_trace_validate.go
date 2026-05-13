package interaction

import (
	"errors"
	"fmt"
)

func ValidateTrace(trace Trace) error {
	// ValidateTrace keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if trace.SchemaVersion != SchemaVersion {
		return errors.New("interaction trace has unsupported schema_version")
	}
	if err := validateTraceHeader(trace); err != nil {
		return err
	}
	if err := validateTraceEvents(trace); err != nil {
		return err
	}
	return validateOrdering(trace.Events)
}

func validateTraceHeader(trace Trace) error {
	// validateTraceHeader keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if err := validateSafeID("task_id", trace.TaskID); err != nil {
		return err
	}
	if !validSourceType(trace.SourceType) {
		return fmt.Errorf("unsupported source_type %q", trace.SourceType)
	}
	if len(trace.Events) == 0 {
		return errors.New("interaction trace requires events")
	}
	return nil
}

func validateTraceEvents(trace Trace) error {
	// validateTraceEvents keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	for _, event := range trace.Events {
		if event.TaskID != trace.TaskID {
			return errors.New("interaction trace event task_id mismatch")
		}
		if err := ValidateEvent(event); err != nil {
			return err
		}
	}
	return nil
}
