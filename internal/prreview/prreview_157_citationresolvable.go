package prreview

func citationResolvable(packet Packet, citation Citation) bool {
	// citationResolvable keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if !citationHasAnchor(citation) {
		return false
	}
	if resolvable, ok := citationRefResolvable(packet, citation); ok {
		return resolvable
	}
	return citation.SourceDigest != ""
}
