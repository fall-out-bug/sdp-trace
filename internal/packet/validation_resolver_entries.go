package packet

import "strings"

func (v *bundleValidator) indexResolverEntries() {
	for _, resolver := range v.bundle.Manifest.Resolvers {
		v.indexResolverEntry(resolver)
	}
}

func (v *bundleValidator) indexResolverEntry(resolver ResolverEntry) {
	if strings.TrimSpace(resolver.Ref) == "" {
		return
	}

	v.resolverByRef[resolver.Ref] = resolver.Resolver
}
