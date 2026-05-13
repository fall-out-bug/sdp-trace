package prreview

func citationHasAnchor(citation Citation) bool {
	return citation.ContextRefID != "" || citation.SourceDigest != ""
}
