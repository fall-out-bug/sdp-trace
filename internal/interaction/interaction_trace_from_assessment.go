package interaction

import (
	"time"
)

func traceFromAssessment(taskID, sourceType string, events []Event, now time.Time, state *traceAssessment) Trace {
	// traceFromAssessment keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	stamp := now.UTC().Format(time.RFC3339)
	return Trace{
		SchemaVersion:     SchemaVersion,
		TraceID:           "it-" + randomHex(12),
		TaskID:            taskID,
		SourceType:        sourceType,
		CompletenessState: state.completeness,
		Events:            events,
		AssessmentState:   state.assessment,
		NotAssessed:       state.notAssessed,
		CannotVerify:      state.cannotVerify,
		CreatedAt:         stamp,
		UpdatedAt:         stamp,
	}
}
