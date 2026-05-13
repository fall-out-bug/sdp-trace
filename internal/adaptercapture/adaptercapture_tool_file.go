package adaptercapture

func toolDepthCondition(run RunEvidence) Condition {
	// Tool-depth evidence records whether adapter events captured nested tool use
	// without turning absent depth into a pass.
	if hasRequired(run, "tool_call") && !hasEvent(run.AdapterEvents, "tool_call") {
		if unsupported(run, "tool_call") {

			return Condition{ID: "tool_call_depth_visible", State: StateUnsupported, ReasonCode: "tool_event_unsupported", Reason: "adapter declares no tool-call capability", NextAction: "Use an adapter with tool-call capture capability."}
		}

		return Condition{ID: "tool_call_depth_visible", State: StateMissingTelemetry, ReasonCode: "tool_event_missing", Reason: "required tool-call adapter event is missing", NextAction: "Capture tool_call adapter events or mark the observer unsupported."}
	}
	return pass("tool_call_depth_visible", "tool_call_depth_visible", "required tool-call families are captured or not required")
}

func fileMutationCondition(run RunEvidence) Condition {
	// File mutation evidence keeps observed changes distinct from unsupported or
	// missing mutation telemetry.
	for _, event := range run.AdapterEvents {
		if fileMutationCorrelationMissing(event) {

			return cannotVerify("file_mutation_correlated", "file_mutation_source_missing", "file mutation is not correlated with source baseline and run id", "Record source baseline and run id correlation for file mutation events.")
		}
	}
	return pass("file_mutation_correlated", "file_mutation_correlated", "file mutation evidence is correlated with source baseline and run id")
}

func fileMutationCorrelationMissing(event AdapterEvent) bool {
	return event.EventType == "file_mutation" && (event.SourceBaseline == "" || event.RunID == "")
}
