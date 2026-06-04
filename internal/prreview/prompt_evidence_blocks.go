package prreview

import (
	"encoding/json"
	"fmt"
	"strings"
)

func appendPromptPacketJSON(b *strings.Builder, packet Packet) error {
	// The packet header is included verbatim so reviewers can echo the digest,
	// commits, and CI state without depending on out-of-band workflow context.
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
	// Ref order follows packet semantics: diff first, then metadata, context,
	// and verification. A single unreadable ref invalidates the prompt.
	for _, ref := range refs {
		if err := appendPromptSafeRef(b, packetDir, ref.label, ref.ref); err != nil {
			return err
		}
	}
	return nil
}

func appendPromptSafeRef(b *strings.Builder, packetDir, label string, ref SafeRef) error {
	// Prompt text is assembled only after the ref path and digest both pass.
	// The rendered block carries the safe ref id for reviewer citations.
	data, err := readPacketRef(packetDir, ref)
	if err != nil {
		return err
	}
	fmt.Fprintf(b, "\n%s ref %s (%s):\n```%s\n%s\n```\n", label, ref.ID, ref.Ref, promptFenceType(ref), string(data))
	return nil
}

func promptFenceType(ref SafeRef) string {
	// Markdown fences are presentation hints only; evidence identity still comes
	// from the packet ref id and SHA-256 digest.
	switch ref.ContentType {
	case ContentUnifiedDiff:
		return "diff"
	case ContentJSON:
		return "json"
	default:
		return "text"
	}
}
