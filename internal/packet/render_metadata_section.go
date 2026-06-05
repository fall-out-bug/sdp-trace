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

func renderMetadata(out *bytes.Buffer, packet Packet) {
	// renderMetadata keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	fmt.Fprintf(out, "## Packet Metadata\n\n")
	fmt.Fprintf(out, "| field | value |\n| --- | --- |\n")
	for _, field := range packetMetadataFields(packet) {
		fmt.Fprintf(out, "| %s | %s |\n", field[0], md(field[1]))
	}
	fmt.Fprintln(out)
}

func packetMetadataFields(packet Packet) [][2]string {
	// packetMetadataFields keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	return [][2]string{
		{"packet_id", packet.PacketID},
		{"schema", packet.PacketVersion},
		{"generated_from", packet.SourceChange.URL},
		{"generated_at", packet.GeneratedAt},
		{"authoring_method", packet.AuthoringMethod},
		{"selected_profile", packet.SelectedProfile},
		{"redaction_policy", packet.RedactionPolicy},
		{"bundle_ref", packet.BundleRef},
		{"packet_state", packet.PacketState},
	}
}
