package prreview

func collectPacketRefs(inputDir string, opts PacketOptions) (packetRefs, error) {
	// collectPacketRefs keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	diffRef, err := copyDiffRef(inputDir, opts.DiffPath)
	if err != nil {
		return packetRefs{}, err
	}
	return collectOptionalPacketRefs(inputDir, opts, diffRef)
}
