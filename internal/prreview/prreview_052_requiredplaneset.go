package prreview

func requiredPlaneSet(planes []string) map[string]bool {
	// requiredPlaneSet keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	required := map[string]bool{}
	for _, plane := range planes {
		if plane != "" {

			required[plane] = true
		}
	}
	return required
}
