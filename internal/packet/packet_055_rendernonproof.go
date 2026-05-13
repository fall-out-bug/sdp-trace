package packet

import (
	"bytes"

	"fmt"
	"strings"
)

func renderNonProof(out *bytes.Buffer, packet Packet) {
	// renderNonProof keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	fmt.Fprintf(out, "## What This Packet Does Not Prove\n\n")
	if strings.TrimSpace(packet.NonApproval) != "" {
		fmt.Fprintf(out, "%s\n\n", packet.NonApproval)
		return
	}
	fmt.Fprintf(out, "This packet does not approve merge, release, compliance, production trust, semantic correctness, or signed external trust.\n\n")
}
