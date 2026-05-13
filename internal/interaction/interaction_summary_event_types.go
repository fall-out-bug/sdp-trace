package interaction

func summarizeTraceEventTypeCounters(eventType string, summary *traceSummaryCounter) {
	// summarizeTraceEventTypeCounters keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	switch eventType {
	case "plan_rejected":
		summary.planRejectionCount++
	case "clarification_request":
		summary.clarificationCount++
	case "clarification_answer":
		summary.clarificationCount++
	}
}

func summarizeTraceReferenceAndCorrection(event Event, summary *traceSummaryCounter) {
	// summarizeTraceReferenceAndCorrection keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if eventIsUnreferenced(event) {
		summary.unreferencedCount++
	}
	if summary.assignmentObserved && isPostAssignmentCorrection(event.EventType) {
		summary.correctionsAfterAssignment++
	}
}

func eventIsUnreferenced(event Event) bool {
	return len(event.ReferenceRefs) == 0 || event.State == StateUnreferenced
}

func isPostAssignmentCorrection(eventType string) bool {
	// isPostAssignmentCorrection keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	switch eventType {
	case "corrective_feedback", "boundary_violation", "evidence_correction":
		return true
	default:
		return false
	}
}
