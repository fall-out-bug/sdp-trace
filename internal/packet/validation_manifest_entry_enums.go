package packet

func (v *bundleValidator) validateManifestEntryEnums(entry BundleEntry) {
	if !retainedForms[entry.RetainedForm] {
		v.add("manifest entry %q has unknown retained_form %q", entry.Ref, entry.RetainedForm)
	}
	if !redactionStatuses[entry.RedactionStatus] {
		v.add("manifest entry %q has unknown redaction_status %q", entry.Ref, entry.RedactionStatus)
	}
}
