package packet

func (v *bundleValidator) validateManifestEntryEnums(entry BundleEntry) {
	// validateManifestEntryEnums keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	if !retainedForms[entry.RetainedForm] {
		v.add("manifest entry %q has unknown retained_form %q", entry.Ref, entry.RetainedForm)
	}
	if !redactionStatuses[entry.RedactionStatus] {
		v.add("manifest entry %q has unknown redaction_status %q", entry.Ref, entry.RedactionStatus)
	}
}
