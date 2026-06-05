package packet

func (v *bundleValidator) indexManifest() {
	v.indexManifestEntries()
	v.indexResolverEntries()
}

func (v *bundleValidator) indexManifestEntries() {
	for _, entry := range v.bundle.Manifest.Entries {
		if v.indexManifestEntry(entry) {
			v.indexValidManifestEntry(entry)
		}
	}
}
