package prreview

// Review status helpers convert raw reviewer outcomes into trust coverage.
//
// A positive reviewer status is not enough by itself: coverage is usable only
// when retained raw-output evidence is present. Missing evidence degrades the
// plane to cannot_verify with exact repair instructions.
func planeResult(result ReviewerResult) PlaneResult {
	pr := PlaneResult{Plane: result.Plane, Status: result.Status, RunID: result.ReviewRunID}
	if reviewerResultUsable(result) {
		pr.Usable = true
		return pr
	}
	if reviewerStatusUsable(result.Status) {
		pr.Status = StatusCannotVerify
		pr.Reason, pr.NextAction = reviewerStatusAction(result.Status)
		return pr
	}
	pr.Reason, pr.NextAction = reviewerStatusAction(result.Status)
	return pr
}

func reviewerStatusUsable(status string) bool {
	return status == StatusFindingsReported || status == StatusNoFindings
}

// Usable review coverage requires both a positive reviewer status and retained
// raw output evidence. A hand-authored status without a digest-bound output
// reference stays unverifiable.
func reviewerResultUsable(result ReviewerResult) bool {
	return reviewerStatusUsable(result.Status) && result.RawOutputRef != nil
}

func reviewerStatusAction(status string) (string, string) {
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
