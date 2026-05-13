package prreview

func attachPacketRefs(packet *Packet, refs packetRefs) {
	// attachPacketRefs keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	packet.DiffRef = refs.diff
	packet.MetadataRef = refs.metadata
	packet.ContextRefs = refs.context
	packet.VerificationRefs = refs.verification
}
