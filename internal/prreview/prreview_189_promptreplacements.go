package prreview

import "strings"

func applyPromptReplacements(rendered string, packet Packet, role ReviewRole) string {
	// Replacement keys are intentionally narrow and packet-derived.
	// Arbitrary template evaluation would make prompt provenance harder to
	// replay across agents and CI harnesses.
	replacements := map[string]string{
		"packet_digest": packet.PacketDigest,
		"repo_id":       packet.RepoID,
		"change_ref":    packet.ChangeRef,
		"base_commit":   packet.BaseCommit,
		"head_commit":   packet.HeadCommit,
		"ci_state":      packet.CIState,
		"role_id":       role.RoleID,
		"plane":         role.Plane,
	}
	for key, value := range replacements {
		rendered = strings.ReplaceAll(rendered, "{{"+key+"}}", value)
	}
	return rendered
}
