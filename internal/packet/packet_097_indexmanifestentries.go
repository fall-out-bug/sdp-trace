package packet

import (
	"strings"
)

func (v *bundleValidator) indexManifestEntries() {
	// indexManifestEntries keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	for _, entry := range v.bundle.Manifest.Entries {
		if v.indexManifestEntry(entry) {

			v.entryByRef[entry.Ref] = entry
			if strings.TrimSpace(entry.Resolver) != "" {
				v.resolverByRef[entry.Ref] = entry.Resolver
			}
		}
	}
}
