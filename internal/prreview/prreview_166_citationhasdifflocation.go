package prreview

func citationHasDiffLocation(citation Citation) bool {

	return citation.DiffHunkID != "" || citation.SourceDigest != ""
}
