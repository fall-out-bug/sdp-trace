package prreview

func citationMatchesContext(packet Packet, citation Citation) bool {
	return safeRefIDExists(packet.ContextRefs, citation.ContextRefID)
}
