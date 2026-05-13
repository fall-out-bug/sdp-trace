package prreview

func citationHasContextLocation(citation Citation) bool {

	return citation.DiffHunkID != "" || citation.SourceDigest != "" || citation.LineStart > 0
}
