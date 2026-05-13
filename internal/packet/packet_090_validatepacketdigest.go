package packet

import (
	"strings"
)

func (v *bundleValidator) validatePacketDigest() {
	// validatePacketDigest keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if strings.TrimSpace(v.bundle.Manifest.PacketDigest) == "" {
		v.add("manifest.packet_digest is required")
	} else if digest := PacketDigest(v.bundle.Packet); digest != "" && v.bundle.Manifest.PacketDigest != digest {

		v.add("manifest.packet_digest does not match packet content")
	}
}
