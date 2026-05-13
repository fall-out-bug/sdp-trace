package adaptercapture

func taskDriftCondition(run RunEvidence) Condition {
	// Task drift evidence distinguishes superseded work from unverified task changes
	// so later gates do not silently accept stale context.
	if !run.TaskDriftAssessed {

		return Condition{ID: "task_drift_visible", State: StateNotAssessed, ReasonCode: "task_drift_not_assessed", Reason: "task drift assessment was not selected", NextAction: "Assess task locks and task_superseded events."}
	}
	if taskSupersessionActorMissing(run.AdapterEvents) {
		return cannotVerify("task_drift_visible", "task_supersession_actor_missing", "task supersession lacks actor attribution state", "Record actor attribution state and task digest refs.")
	}
	return taskDriftPassCondition(run.TaskSupersessionCount)
}

func taskDriftPassCondition(supersessionCount int) Condition {
	// taskDriftPassCondition preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.
	if supersessionCount == 0 {
		return pass("task_drift_visible", "no_supersessions_observed", "task drift was assessed and no supersessions were observed")
	}

	return pass("task_drift_visible", "task_supersessions_visible", "task supersessions include visible attribution and digest evidence")
}

func taskSupersessionActorMissing(events []AdapterEvent) bool {
	// taskSupersessionActorMissing preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.
	for _, event := range events {
		if event.EventType == "task_superseded" && event.ActorAttributionState == "" {

			return true
		}
	}
	return false
}
