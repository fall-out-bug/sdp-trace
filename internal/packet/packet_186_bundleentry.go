package packet

func bundleEntry(ref, sourceClass, resolver, retainedForm string) BundleEntry {
	// bundleEntry keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	resolver = redactSecretLike(resolver)

	return BundleEntry{
		Ref:             ref,
		SourceClass:     sourceClass,
		Digest:          digestPlaceholder(ref + resolver),
		RetainedForm:    retainedForm,
		RedactionStatus: "not_needed",
		Resolver:        resolver,
		ArtifactAccess:  "present",
	}
}
