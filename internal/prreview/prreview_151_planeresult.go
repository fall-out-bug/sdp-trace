package prreview

func planeResult(result ReviewerResult) PlaneResult {
	// planeResult keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

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
