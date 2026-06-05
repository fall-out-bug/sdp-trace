package prreview

import (
	"path/filepath"
)

func BuildPacket(opts PacketOptions) (Packet, error) {
	// BuildPacket keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if err := validatePacketOptions(opts); err != nil {
		return Packet{}, err
	}
	if err := ensureNewDir(opts.OutDir); err != nil {
		return Packet{}, err
	}
	return buildPacketInPreparedDir(opts)
}

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
