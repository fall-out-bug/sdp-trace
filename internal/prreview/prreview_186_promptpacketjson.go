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
