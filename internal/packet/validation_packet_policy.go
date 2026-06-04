package packet

func (v *bundleValidator) validatePacketPolicyMetadata() {
	v.requireNonEmpty(v.bundle.Packet.NonApproval, "packet.non_approval is required")
	v.requireKnown(packetStates, v.bundle.Packet.PacketState, "packet.packet_state has unknown value %q")
	v.requireKnown(authoringMethods, v.bundle.Packet.AuthoringMethod, "packet.authoring_method has unknown value %q")
	v.validateProjectionMetadata()
}

func (v *bundleValidator) requireKnown(known map[string]bool, value string, format string) {
	if !known[value] {
		v.add(format, value)
	}
}
