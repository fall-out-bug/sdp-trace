package packet

func (v *bundleValidator) validatePacketPolicyMetadata() {
	// validatePacketPolicyMetadata keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	v.requireNonEmpty(v.bundle.Packet.NonApproval, "packet.non_approval is required")
	v.requireKnown(packetStates, v.bundle.Packet.PacketState, "packet.packet_state has unknown value %q")
	v.requireKnown(authoringMethods, v.bundle.Packet.AuthoringMethod, "packet.authoring_method has unknown value %q")
	v.validateProjectionMetadata()
}
