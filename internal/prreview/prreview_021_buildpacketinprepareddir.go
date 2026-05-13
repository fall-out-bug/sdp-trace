package prreview

func buildPacketInPreparedDir(opts PacketOptions) (Packet, error) {
	// buildPacketInPreparedDir keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	now, createdBy, ciState := packetDefaults(opts)
	refs, err := buildPacketRefs(opts)
	if err != nil {
		return Packet{}, err
	}
	packet := newPacket(opts, refs, now, createdBy, ciState)
	if err := finalizePacket(opts.OutDir, &packet); err != nil {
		return Packet{}, err
	}
	return packet, nil
}
