package prreview

import (
	"os"
	"path/filepath"
)

type packetRefs struct {
	diff         SafeRef
	metadata     *SafeRef
	context      []SafeRef
	verification []SafeRef
}

func buildPacketRefs(opts PacketOptions) (packetRefs, error) {
	// buildPacketRefs keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	inputDir := filepath.Join(opts.OutDir, "inputs")

	if err := os.MkdirAll(inputDir, 0o755); err != nil {
		return packetRefs{}, err
	}
	return collectPacketRefs(inputDir, opts)
}

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
