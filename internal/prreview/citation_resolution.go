package prreview

// Citation resolution requires an anchor and then resolves known packet refs.
//
// Unknown refs are accepted only when the reviewer also provides a source
// digest, preserving the distinction between anchored evidence and a bare label.
func citationResolvable(packet Packet, citation Citation) bool {
	if !citationHasAnchor(citation) {
		return false
	}
	if resolvable, ok := citationRefResolvable(packet, citation); ok {
		return resolvable
	}
	return citation.SourceDigest != ""
}

func citationHasAnchor(citation Citation) bool {
	return citation.ContextRefID != "" || citation.SourceDigest != ""
}

func citationRefResolvable(packet Packet, citation Citation) (bool, bool) {
	for _, resolver := range citationResolvers {
		if resolver.matches(packet, citation) {
			return resolver.resolvable(citation), true
		}
	}
	return false, false
}

type citationResolver struct {
	matches    func(Packet, Citation) bool
	resolvable func(Citation) bool
}

var citationResolvers = []citationResolver{
	{matches: citationMatchesDiff, resolvable: citationHasDiffLocation},
	{matches: citationMatchesContext, resolvable: citationHasContextLocation},
	{matches: citationMatchesVerification, resolvable: citationHasVerificationLocation},
}
