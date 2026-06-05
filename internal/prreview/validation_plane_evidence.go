package prreview

// planeResultWithModelCheck projects one reviewer result into plane coverage
// while enforcing the declared model identity from the review profile.
func planeResultWithModelCheck(role ReviewRole, result ReviewerResult) PlaneResult {
	best := planeResult(result)
	if best.Usable && modelMismatchWithoutFallback(role, result) {
		best.Usable = false
		best.Status = StatusCannotVerify
		best.Reason = "model_identity_mismatch"
		best.NextAction = "Rerun the reviewer or record fallback provenance for the observed model."
	}
	return best
}

// planeCannotVerify identifies failed reviewer states that are stronger than
// not_assessed because the harness observed an attempted but unreplayable run.
func planeCannotVerify(status string) bool {
	switch status {
	case StatusCannotVerify, StatusTimedOut, StatusEmptyOutput, StatusOffTask, StatusParseFailed:
		return true
	default:
		return false
	}
}
