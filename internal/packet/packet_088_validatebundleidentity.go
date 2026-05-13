package packet

func (v *bundleValidator) validateBundleIdentity() {
	// validateBundleIdentity keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	v.requireNonEmpty(v.bundle.Packet.PacketID, "packet.packet_id is required")
	v.requireNonEmpty(v.bundle.Packet.BundleRef, "packet.bundle_ref is required")
	v.requireNonEmpty(v.bundle.Manifest.BundleID, "manifest.bundle_id is required")

	if v.bundle.Packet.BundleRef != "" && v.bundle.Manifest.BundleID != "" && v.bundle.Packet.BundleRef != v.bundle.Manifest.BundleID {
		v.add("packet.bundle_ref %q must match manifest.bundle_id %q", v.bundle.Packet.BundleRef, v.bundle.Manifest.BundleID)
	}
}
