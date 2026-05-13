package packet

func packetMetadataFields(packet Packet) [][2]string {
	// packetMetadataFields keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	return [][2]string{
		{"packet_id", packet.PacketID},
		{"schema", packet.PacketVersion},
		{"generated_from", packet.SourceChange.URL},
		{"generated_at", packet.GeneratedAt},
		{"authoring_method", packet.AuthoringMethod},
		{"selected_profile", packet.SelectedProfile},
		{"redaction_policy", packet.RedactionPolicy},
		{"bundle_ref", packet.BundleRef},
		{"packet_state", packet.PacketState},
	}
}
