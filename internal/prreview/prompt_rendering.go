package prreview

import (
	"fmt"
	"os"
	"strings"
)

func renderPrompt(packet Packet, role ReviewRole, packetDir string) (string, error) {
	// Reviewer prompts carry two evidence classes: the selected template and
	// immutable packet refs. This keeps model review food bound to the packet.
	// If packet evidence cannot be replayed, the caller records cannot_verify.
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
	// Template absence is allowed for imported/manual evidence paths.
	// Template read failure is not downgraded because the configured runner
	// would otherwise review a different prompt from the declared profile.
	if strings.TrimSpace(role.PromptTemplateRef) == "" {
		return "", nil
	}
	data, err := os.ReadFile(role.PromptTemplateRef)
	if err != nil {
		return "", fmt.Errorf("%w: prompt_template_ref", errPromptTemplateCannotVerify)
	}
	return applyPromptReplacements(string(data), packet, role), nil
}

func renderPromptEvidence(packet Packet, packetDir string) (string, error) {
	// PacketDir is optional for legacy imported runs; when present, every ref
	// is read through digest-checked safe paths before it enters the prompt.
	// This prevents reviewers from seeing stale or path-traversed evidence.
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
