package prreview

func packetDigest(packet Packet) (string, error) {
	// packetDigest keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	canonical := packet
	canonical.PacketDigest = ""
	return digestJSON(canonical)
}
