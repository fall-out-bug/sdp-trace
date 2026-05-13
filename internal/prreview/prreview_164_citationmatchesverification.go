package prreview

func citationMatchesVerification(packet Packet, citation Citation) bool {
	return safeRefIDExists(packet.VerificationRefs, citation.ContextRefID)
}
