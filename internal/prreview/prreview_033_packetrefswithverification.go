package prreview

func packetRefsWithVerification(inputDir string, opts PacketOptions, diffRef SafeRef, metadataRef *SafeRef, contextRefs []SafeRef) (packetRefs, error) {
	// packetRefsWithVerification keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	verificationRefs, err := packetVerificationRefs(inputDir, opts.VerificationPaths)
	if err != nil {
		return packetRefs{}, err
	}
	return packetRefs{diff: diffRef, metadata: metadataRef, context: contextRefs, verification: verificationRefs}, nil
}
