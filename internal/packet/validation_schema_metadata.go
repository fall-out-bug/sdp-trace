package packet

func (v *bundleValidator) validateSchemaMetadata() {
	if v.bundle.Packet.PacketVersion != PacketSchemaVersion {
		v.add("packet.packet_version must be %q", PacketSchemaVersion)
	}

	if v.bundle.Manifest.SchemaVersion != BundleSchemaVersion {
		v.add("manifest.schema_version must be %q", BundleSchemaVersion)
	}
}
