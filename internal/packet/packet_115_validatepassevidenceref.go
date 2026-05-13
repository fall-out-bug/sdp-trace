package packet

func (v *bundleValidator) validatePassEvidenceRef(rowID, ref string, entry BundleEntry) {
	// validatePassEvidenceRef keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	if entryExpired(entry, v.now) {
		v.add("%s pass cites expired artifact ref %q", rowID, ref)
	}
	if passRefUnverifiable(entry) {
		v.add("%s pass cites unverifiable artifact ref %q", rowID, ref)
	}
}
