package prreview

func planeResultWithModelCheck(role ReviewRole, result ReviewerResult) PlaneResult {
	// planeResultWithModelCheck keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	best := planeResult(result)
	if best.Usable && modelMismatchWithoutFallback(role, result) {
		best.Usable = false
		best.Status = StatusCannotVerify
		best.Reason = "model_identity_mismatch"
		best.NextAction = "Rerun the reviewer or record fallback provenance for the observed model."
	}
	return best
}
