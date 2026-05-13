package prreview

func collectOptionalPacketRefs(inputDir string, opts PacketOptions, diffRef SafeRef) (packetRefs, error) {
	// collectOptionalPacketRefs keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	metadataRef, err := optionalMetadataRef(inputDir, opts.MetadataPath)
	if err != nil {
		return packetRefs{}, err
	}
	contextRefs, err := packetContextRefs(inputDir, opts.ContextPaths)
	if err != nil {
		return packetRefs{}, err
	}
	return packetRefsWithVerification(inputDir, opts, diffRef, metadataRef, contextRefs)
}
