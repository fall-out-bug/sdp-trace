package prreview

import (
	"os"

	"path/filepath"
)

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
