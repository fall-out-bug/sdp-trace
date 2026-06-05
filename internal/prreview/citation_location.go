package prreview

// Citation location checks are intentionally field-specific.
//
// Diff citations require a hunk or digest. Context citations may point to a
// hunk, digest, or line. Verification citations require a line or digest; a
// diff hunk alone does not anchor verifier output.
func citationHasDiffLocation(citation Citation) bool {
	return citation.DiffHunkID != "" || citation.SourceDigest != ""
}

func citationHasContextLocation(citation Citation) bool {
	return citation.DiffHunkID != "" || citation.SourceDigest != "" || citation.LineStart > 0
}

func citationHasVerificationLocation(citation Citation) bool {
	return citation.SourceDigest != "" || citation.LineStart > 0
}
