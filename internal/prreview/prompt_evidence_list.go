package prreview

type promptEvidenceRef struct {
	label string
	ref   SafeRef
}

func promptEvidenceRefs(packet Packet) []promptEvidenceRef {
	// Diff evidence is mandatory in packet construction; optional families are
	// appended only when the packet declares them.
	refs := []promptEvidenceRef{{label: "diff", ref: packet.DiffRef}}
	refs = appendOptionalMetadataRef(refs, packet.MetadataRef)
	refs = appendPromptRefs(refs, "context", packet.ContextRefs)
	return appendPromptRefs(refs, "verification", packet.VerificationRefs)
}

func appendOptionalMetadataRef(refs []promptEvidenceRef, ref *SafeRef) []promptEvidenceRef {
	// Metadata is represented as a pointer in the packet, so nil means the
	// missing field was already captured as not_assessed during packet build.
	if ref == nil {
		return refs
	}
	return append(refs, promptEvidenceRef{label: "metadata", ref: *ref})
}

func appendPromptRefs(refs []promptEvidenceRef, label string, safeRefs []SafeRef) []promptEvidenceRef {
	// Context and verification refs preserve packet order so prompt replay
	// matches the evidence bundle exactly.
	for _, ref := range safeRefs {
		refs = append(refs, promptEvidenceRef{label: label, ref: ref})
	}
	return refs
}
