package prreview

func reviewerStatusAction(status string) (string, string) {
	// reviewerStatusAction keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	actions := map[string][2]string{
		StatusNotAssessed: {"reviewer_not_assessed", "Run a configured reviewer or import a usable result for this plane."},
		StatusTimedOut:    {"reviewer_timed_out", "Increase timeout or replace the reviewer for this plane."},
		StatusEmptyOutput: {"reviewer_empty_output", "Retry with a shorter bounded prompt or replace the reviewer."},
		StatusOffTask:     {"reviewer_off_task", "Rerun with the frozen packet and required output schema."},
		StatusParseFailed: {"reviewer_parse_failed", "Rerun with JSON-only output matching the required schema."},
	}
	if action, ok := actions[status]; ok {
		return action[0], action[1]
	}
	if reviewerStatusUsable(status) {
		return "reviewer_output_not_retained", "Attach digest-bound reviewer output before counting this plane."
	}
	return "reviewer_cannot_verify", "Replace or rerun the reviewer."
}
