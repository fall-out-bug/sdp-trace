package prreview

func citationMatchesDiff(packet Packet, citation Citation) bool {
	return citation.ContextRefID == packet.DiffRef.ID || citation.ContextRefID == "diff"
}
