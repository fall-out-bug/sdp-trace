package prreview

func preview(packet Packet, profile ReviewProfile) *RunPreview {
	// preview keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	roles := make([]PreviewRole, 0, len(profile.Roles))
	for _, role := range profile.Roles {
		roles = append(roles, previewRole(role))
	}
	return &RunPreview{SchemaVersion: SchemaVersionRunSet, PacketDigest: packet.PacketDigest, Roles: roles}
}
