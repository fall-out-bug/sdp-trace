package prreview

func reviewerOutputMismatched(parsed ReviewerResult, role ReviewRole, packet Packet) bool {

	return parsed.PacketDigest != packet.PacketDigest || parsed.Plane != role.Plane || parsed.RoleID != role.RoleID
}
