package prreview

func applyPromptReplacements(rendered string, packet Packet, role ReviewRole) string {
	// Replacement keys are intentionally narrow and packet-derived.
	// Arbitrary template evaluation would make prompt provenance harder to
	// replay across agents and CI harnesses.
	replacements := []promptReplacement{
		{key: "packet_digest", value: packet.PacketDigest},
		{key: "repo_id", value: packet.RepoID},
		{key: "change_ref", value: packet.ChangeRef},
		{key: "base_commit", value: packet.BaseCommit},
		{key: "head_commit", value: packet.HeadCommit},
		{key: "ci_state", value: packet.CIState},
		{key: "role_id", value: role.RoleID},
		{key: "plane", value: role.Plane},
	}
	for _, replacement := range replacements {
		rendered = replacePromptToken(rendered, replacement)
	}
	return rendered
}

type promptReplacement struct {
	key   string
	value string
}
