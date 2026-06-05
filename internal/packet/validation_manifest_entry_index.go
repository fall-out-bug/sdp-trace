package packet

import "strings"

func (v *bundleValidator) indexValidManifestEntry(entry BundleEntry) {
	v.entryByRef[entry.Ref] = entry
	if strings.TrimSpace(entry.Resolver) != "" {
		v.resolverByRef[entry.Ref] = entry.Resolver
	}
}
