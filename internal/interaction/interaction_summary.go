package interaction

func SummarizeTrace(trace Trace) Summary {
	// SummarizeTrace keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	summary := newTraceSummaryCounter()
	for _, event := range trace.Events {
		summarizeTraceEvent(event, summary)
	}
	notAssessed := append([]string{}, trace.NotAssessed...)
	if !summary.assignmentObserved {
		notAssessed = append(notAssessed, "task_assignment event absent; post-assignment correction count is not assessed")
	}

	return traceSummary(trace, summary, notAssessed)
}
func traceSummary(trace Trace, summary *traceSummaryCounter, notAssessed []string) Summary {
	// traceSummary keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	counts := traceSummaryCounts(trace, summary)
	return Summary{
		SchemaVersion:   SchemaVersion,
		TaskID:          trace.TaskID,
		TraceID:         trace.TraceID,
		AssessmentState: trace.AssessmentState,

		EventCount:             counts.events,
		FrictionCounts:         counts.friction,
		CorrectionAfterTask:    counts.corrections,
		PlanRejectionCount:     counts.planRejections,
		ClarificationTurnCount: counts.clarifications,
		UnreferencedEventCount: counts.unreferenced,

		NotAssessed:  notAssessed,
		CannotVerify: trace.CannotVerify,
	}
}
