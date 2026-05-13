package packet

import (
	"strings"
)

func (v *bundleValidator) indexManifestEntry(entry BundleEntry) bool {
	// indexManifestEntry keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if strings.TrimSpace(entry.Ref) == "" {
		v.add("manifest entry has empty ref")
		return false
	}

	v.validateManifestEntryEnums(entry)
	return true
}
