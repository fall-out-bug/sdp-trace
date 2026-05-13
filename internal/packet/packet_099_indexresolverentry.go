package packet

import (
	"strings"
)

func (v *bundleValidator) indexResolverEntry(resolver ResolverEntry) {
	// indexResolverEntry keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if strings.TrimSpace(resolver.Ref) == "" {
		return
	}

	v.resolverByRef[resolver.Ref] = resolver.Resolver
}
