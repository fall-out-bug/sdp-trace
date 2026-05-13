package packet

func (v *bundleValidator) validateSchemaMetadata() {
	// validateSchemaMetadata keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if v.bundle.Packet.PacketVersion != PacketSchemaVersion {
		v.add("packet.packet_version must be %q", PacketSchemaVersion)
	}

	if v.bundle.Manifest.SchemaVersion != BundleSchemaVersion {
		v.add("manifest.schema_version must be %q", BundleSchemaVersion)
	}
}
