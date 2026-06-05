package packet

import "strings"

func (v *bundleValidator) validateBundleIdentity() {
	v.requireNonEmpty(v.bundle.Packet.PacketID, "packet.packet_id is required")
	v.requireNonEmpty(v.bundle.Packet.BundleRef, "packet.bundle_ref is required")
	v.requireNonEmpty(v.bundle.Manifest.BundleID, "manifest.bundle_id is required")

	if v.bundle.Packet.BundleRef != "" && v.bundle.Manifest.BundleID != "" && v.bundle.Packet.BundleRef != v.bundle.Manifest.BundleID {
		v.add("packet.bundle_ref %q must match manifest.bundle_id %q", v.bundle.Packet.BundleRef, v.bundle.Manifest.BundleID)
	}
}

func (v *bundleValidator) requireNonEmpty(value string, message string) {
	if strings.TrimSpace(value) == "" {
		v.add("%s", message)
	}
}
