package prreview

func requiredPlaneResult(plane string, roleByID map[string]ReviewRole, runs RunSet, reasons, nextActions *[]string) PlaneResult {
	// requiredPlaneResult keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	result := bestPlaneResult(plane, roleByID, runs)
	appendPlaneValidationNotes(result, reasons, nextActions)
	return result
}
