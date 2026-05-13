package prreview

func packetDigestMismatch(packet Packet, runs RunSet, ledger Ledger) bool {
	return runs.PacketDigest != packet.PacketDigest || ledger.PacketDigest != packet.PacketDigest
}
