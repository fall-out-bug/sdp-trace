package packet

func (v *bundleValidator) indexManifest() {
	v.indexManifestEntries()
	v.indexResolverEntries()
}
