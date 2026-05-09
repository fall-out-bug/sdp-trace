package prreview

import (
	"os"
	"strings"
)

func renderPrompt(packet Packet, role ReviewRole) (string, error) {
	if strings.TrimSpace(role.PromptTemplateRef) == "" {
		return "", nil
	}
	data, err := os.ReadFile(role.PromptTemplateRef)
	if err != nil {
		return "", err
	}
	rendered := string(data)
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
	return rendered, nil
}
