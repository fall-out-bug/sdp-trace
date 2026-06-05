package prreview

// Citation ref matching binds reviewer citations to packet-local evidence.
//
// The diff ref also accepts the stable `diff` alias used in reviewer output,
// while context and verification citations must match copied packet refs by ID.
func citationMatchesDiff(packet Packet, citation Citation) bool {
	return citation.ContextRefID == packet.DiffRef.ID || citation.ContextRefID == "diff"
}

func citationMatchesContext(packet Packet, citation Citation) bool {
	return safeRefIDExists(packet.ContextRefs, citation.ContextRefID)
}

func citationMatchesVerification(packet Packet, citation Citation) bool {
	return safeRefIDExists(packet.VerificationRefs, citation.ContextRefID)
}

func safeRefIDExists(refs []SafeRef, id string) bool {
	for _, ref := range refs {
		if id == ref.ID {
			return true
		}
	}
	return false
}
