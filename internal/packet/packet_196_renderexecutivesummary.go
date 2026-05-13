package packet

import (
	"bytes"

	"fmt"
)

func renderExecutiveSummary(out *bytes.Buffer, packet Packet) {
	// renderExecutiveSummary keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	fmt.Fprintf(out, "## Executive Summary\n\n")
	fmt.Fprintf(out, "- Source change: %s %s.\n", packet.SourceChange.Repository, packet.SourceChange.ChangeID)
	fmt.Fprintf(out, "- Packet state: %s.\n", packet.PacketState)
	fmt.Fprintf(out, "- Selected evidence profile: %s.\n", packet.SelectedProfile)
	fmt.Fprintf(out, "- Required rows preserve pass, partial, fail, cannot_verify, not_assessed, and not_in_scope states without a score.\n")
	fmt.Fprintf(out, "- Next decision ownership is recorded separately from approval.\n\n")
}
