package packet

import (
	"bytes"

	"fmt"
)

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
