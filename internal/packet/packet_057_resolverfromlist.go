package packet

func resolverFromList(resolvers []ResolverEntry, ref string) string {
	// resolverFromList keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	for _, resolver := range resolvers {
		if resolver.Ref == ref {

			return resolver.Resolver
		}
	}
	return ""
}
