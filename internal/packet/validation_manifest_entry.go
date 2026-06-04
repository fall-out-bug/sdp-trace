package packet

import "strings"

func (v *bundleValidator) indexManifestEntry(entry BundleEntry) bool {
	if strings.TrimSpace(entry.Ref) == "" {
		v.add("manifest entry has empty ref")
		return false
	}

	v.validateManifestEntryEnums(entry)
	return true
}
