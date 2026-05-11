package prreview

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func renderPrompt(packet Packet, role ReviewRole, packetDir string) (string, error) {
	rendered, err := renderPromptTemplate(packet, role)
	if err != nil {
		return "", err
	}
	evidence, err := renderPromptEvidence(packet, packetDir)
	if err != nil {
		return "", err
	}
	return rendered + evidence, nil
}

func renderPromptTemplate(packet Packet, role ReviewRole) (string, error) {
	if strings.TrimSpace(role.PromptTemplateRef) == "" {
		return "", nil
	}
	data, err := os.ReadFile(role.PromptTemplateRef)
	if err != nil {
		return "", err
	}
	return applyPromptReplacements(string(data), packet, role), nil
}

func applyPromptReplacements(rendered string, packet Packet, role ReviewRole) string {
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

func renderPromptEvidence(packet Packet, packetDir string) (string, error) {
	if strings.TrimSpace(packetDir) == "" {
		return "", nil
	}
	var b strings.Builder
	if err := appendPromptPacketJSON(&b, packet); err != nil {
		return "", err
	}
	if err := appendPromptEvidenceRefs(&b, packetDir, promptEvidenceRefs(packet)); err != nil {
		return "", err
	}
	return b.String(), nil
}

func appendPromptPacketJSON(b *strings.Builder, packet Packet) error {
	b.WriteString("\n\nReview packet JSON:\n```json\n")
	packetJSON, err := json.MarshalIndent(packet, "", "  ")
	if err != nil {
		return fmt.Errorf("%w: packet_json", errPromptEvidenceCannotVerify)
	}
	b.Write(packetJSON)
	b.WriteString("\n```\n")
	return nil
}

func appendPromptEvidenceRefs(b *strings.Builder, packetDir string, refs []promptEvidenceRef) error {
	for _, ref := range refs {
		if err := appendPromptSafeRef(b, packetDir, ref.label, ref.ref); err != nil {
			return err
		}
	}
	return nil
}

type promptEvidenceRef struct {
	label string
	ref   SafeRef
}

func promptEvidenceRefs(packet Packet) []promptEvidenceRef {
	refs := []promptEvidenceRef{{label: "diff", ref: packet.DiffRef}}
	refs = appendOptionalMetadataRef(refs, packet.MetadataRef)
	refs = appendPromptRefs(refs, "context", packet.ContextRefs)
	return appendPromptRefs(refs, "verification", packet.VerificationRefs)
}

func appendOptionalMetadataRef(refs []promptEvidenceRef, ref *SafeRef) []promptEvidenceRef {
	if ref == nil {
		return refs
	}
	return append(refs, promptEvidenceRef{label: "metadata", ref: *ref})
}

func appendPromptRefs(refs []promptEvidenceRef, label string, safeRefs []SafeRef) []promptEvidenceRef {
	for _, ref := range safeRefs {
		refs = append(refs, promptEvidenceRef{label: label, ref: ref})
	}
	return refs
}

func appendPromptSafeRef(b *strings.Builder, packetDir, label string, ref SafeRef) error {
	data, err := readPacketRef(packetDir, ref)
	if err != nil {
		return err
	}
	fmt.Fprintf(b, "\n%s ref %s (%s):\n```%s\n%s\n```\n", label, ref.ID, ref.Ref, promptFenceType(ref), string(data))
	return nil
}

func promptFenceType(ref SafeRef) string {
	switch ref.ContentType {
	case ContentUnifiedDiff:
		return "diff"
	case ContentJSON:
		return "json"
	default:
		return "text"
	}
}

func readPacketRef(packetDir string, ref SafeRef) ([]byte, error) {
	path, err := packetRefPath(packetDir, ref.Ref)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errPromptEvidenceCannotVerify, ref.ID)
	}
	if !digestMatches(data, ref.DigestSHA256) {
		return nil, fmt.Errorf("%w: %s", errPromptEvidenceCannotVerify, ref.ID)
	}
	return data, nil
}

func packetRefPath(packetDir, ref string) (string, error) {
	cleanRef := filepath.Clean(filepath.FromSlash(ref))
	if filepath.IsAbs(cleanRef) || strings.HasPrefix(cleanRef, ".."+string(filepath.Separator)) || cleanRef == ".." {
		return "", fmt.Errorf("%w: unsafe_ref", errPromptEvidenceCannotVerify)
	}
	return filepath.Join(packetDir, cleanRef), nil
}

func digestMatches(data []byte, digest string) bool {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]) == digest
}
