package prreview

func citationRefResolvable(packet Packet, citation Citation) (bool, bool) {
	// citationRefResolvable keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	for _, resolver := range citationResolvers {
		if resolver.matches(packet, citation) {
			return resolver.resolvable(citation), true
		}
	}
	return false, false
}
