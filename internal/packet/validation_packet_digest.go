package packet

import "strings"

func (v *bundleValidator) validatePacketDigest() {
	if strings.TrimSpace(v.bundle.Manifest.PacketDigest) == "" {
		v.add("manifest.packet_digest is required")
	} else if digest := PacketDigest(v.bundle.Packet); digest != "" && v.bundle.Manifest.PacketDigest != digest {
		v.add("manifest.packet_digest does not match packet content")
	}
}
