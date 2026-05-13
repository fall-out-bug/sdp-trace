package prreview

var citationResolvers = []citationResolver{
	{matches: citationMatchesDiff, resolvable: citationHasDiffLocation},
	{matches: citationMatchesContext, resolvable: citationHasContextLocation},
	{matches: citationMatchesVerification, resolvable: citationHasVerificationLocation},
}
