package interaction

func traceSummaryCounts(trace Trace, summary *traceSummaryCounter) traceCounts {
	// traceSummaryCounts keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	return traceCounts{
		events:         len(trace.Events),
		friction:       summary.frictionCounts,
		corrections:    summary.correctionsAfterAssignment,
		planRejections: summary.planRejectionCount,
		clarifications: summary.clarificationCount,
		unreferenced:   summary.unreferencedCount,
	}
}

func newTraceSummaryCounter() *traceSummaryCounter {
	return &traceSummaryCounter{
		frictionCounts: make(map[string]int),
	}
}

func summarizeTraceEvent(event Event, summary *traceSummaryCounter) {
	// summarizeTraceEvent keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	summary.frictionCounts[event.FrictionClass]++
	if event.EventType == "task_assignment" {
		summary.assignmentObserved = true
		return
	}
	summarizeTraceEventTypeCounters(event.EventType, summary)
	summarizeTraceReferenceAndCorrection(event, summary)
}
