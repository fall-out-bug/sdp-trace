package prreview

func citationHasVerificationLocation(citation Citation) bool {

	return citation.SourceDigest != "" || citation.LineStart > 0
}
