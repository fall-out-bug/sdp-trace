package prreview

import (
	"path/filepath"
)

func finalizePacket(outDir string, packet *Packet) error {
	// finalizePacket keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	digest, err := packetDigest(*packet)
	if err != nil {
		return err
	}
	packet.PacketDigest = "sha256:" + digest
	return WriteJSON(filepath.Join(outDir, "packet.json"), *packet)
}
