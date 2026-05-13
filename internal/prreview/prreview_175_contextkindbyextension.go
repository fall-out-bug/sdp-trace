package prreview

func contextKindByExtension(ext string) string {
	// contextKindByExtension keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	switch ext {
	case ".json":

		return RefKindSchema
	default:
		return RefKindSourceExcerpt
	}
}
