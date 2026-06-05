package prreview

// Preview projection is intentionally read-only.
//
// It mirrors the profile role order and exposes only replayable references:
// packet digest, runner command digest, and prompt digest metadata. It must not
// prepare output directories, create run artifacts, or infer reviewer results.
func preview(packet Packet, profile ReviewProfile) *RunPreview {
	roles := make([]PreviewRole, 0, len(profile.Roles))
	for _, role := range profile.Roles {
		roles = append(roles, previewRole(role))
	}
	return &RunPreview{SchemaVersion: SchemaVersionRunSet, PacketDigest: packet.PacketDigest, Roles: roles}
}

func previewRole(role ReviewRole) PreviewRole {
	return PreviewRole{
		RoleID:         role.RoleID,
		Plane:          role.Plane,
		Runner:         role.Runner,
		RequestedModel: defaultString(role.RequestedModel, StateNotAssessed),
		TimeoutSeconds: role.TimeoutSeconds,
		CommandDigest:  commandDigest(role.Command),
		PromptRef:      role.PromptTemplateRef,
		PromptDigest:   promptDigestForRole(role),
	}
}
