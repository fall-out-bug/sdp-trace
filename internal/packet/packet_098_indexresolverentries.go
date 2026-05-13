package packet

func (v *bundleValidator) indexResolverEntries() {
	// indexResolverEntries keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	for _, resolver := range v.bundle.Manifest.Resolvers {

		v.indexResolverEntry(resolver)
	}
}
