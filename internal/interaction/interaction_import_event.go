package interaction

import (
	"errors"
	"fmt"
)

func normalizeTranscriptEvent(event *Event, opts ImportOptions) error {
	// normalizeTranscriptEvent keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if err := validateTranscriptEventTask(*event, opts); err != nil {
		return err
	}
	if err := validateTranscriptEventSource(*event); err != nil {
		return err
	}
	event.Source.SourceType = SourcePreclassifiedTranscript
	if opts.SourceRef != "" {
		event.Source.SourceRef = opts.SourceRef
	}
	return ValidateEvent(*event)
}

func validateTranscriptEventTask(event Event, opts ImportOptions) error {
	// validateTranscriptEventTask keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.
	if event.TaskID != opts.TaskID {

		return fmt.Errorf("event task_id %q does not match import task_id", event.TaskID)
	}
	return nil
}

func validateTranscriptEventSource(event Event) error {
	// validateTranscriptEventSource keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if event.Source.SourceType == SourceAgentReported {
		return errors.New("agent-reported interaction is not accepted as event evidence")
	}
	if event.Source.SourceType != "" && event.Source.SourceType != SourcePreclassifiedTranscript {
		return fmt.Errorf("unsupported source_type %q", event.Source.SourceType)
	}
	return nil
}
