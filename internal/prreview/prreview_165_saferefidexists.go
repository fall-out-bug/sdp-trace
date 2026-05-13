package prreview

func safeRefIDExists(refs []SafeRef, id string) bool {
	// safeRefIDExists keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	for _, ref := range refs {
		if id == ref.ID {

			return true
		}
	}
	return false
}
