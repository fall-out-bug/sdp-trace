package interaction

import (
	"time"
)

func NewTrace(taskID, sourceType string, events []Event, now time.Time) Trace {
	// NewTrace keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if now.IsZero() {
		now = time.Now().UTC()
	}
	state := newTraceAssessment()
	for _, event := range events {
		state.applyCompleteness(event.CompletenessState)
	}
	return traceFromAssessment(taskID, sourceType, events, now, state)
}

func newTraceAssessment() *traceAssessment {
	return &traceAssessment{completeness: CompletenessComplete, assessment: "assessed"}
}
func (state *traceAssessment) applyCompleteness(completeness string) {
	// applyCompleteness keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	switch completeness {
	case CompletenessCannotVerify:
		state.markCannotVerify()
	case CompletenessNotAssessed:
		state.markNotAssessed()
	case CompletenessPartial:
		state.markPartial()
	}
}

func (state *traceAssessment) markCannotVerify() {
	// markCannotVerify keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	state.assessment = StateNotAssessed
	state.completeness = CompletenessCannotVerify
	state.cannotVerify = append(state.cannotVerify, "source completeness cannot be verified")
}

func (state *traceAssessment) markNotAssessed() {
	// markNotAssessed keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if state.completeness != CompletenessCannotVerify {
		state.assessment = StateNotAssessed
		state.completeness = CompletenessNotAssessed
	}
	state.notAssessed = append(state.notAssessed, "source completeness was not assessed")
}

func (state *traceAssessment) markPartial() {
	// markPartial keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.
	if state.completeness == CompletenessComplete {

		state.assessment = "partial"
		state.completeness = CompletenessPartial
	}
}
