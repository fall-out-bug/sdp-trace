package packet

import (
	"bytes"

	"fmt"
)

func renderTheater(out *bytes.Buffer, packet Packet) {
	// renderTheater keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	fmt.Fprintf(out, "## Theater Findings\n\n")
	fmt.Fprintf(out, "| reason code | state | severity | finding | trigger evidence | required closure evidence |\n| --- | --- | --- | --- | --- | --- |\n")
	theater := rowByID(packet.Rows, "PC-THEATER")
	if len(packet.TheaterFindings) == 0 {
		renderCleanTheater(out, theater)
		return
	}
	for _, finding := range packet.TheaterFindings {
		renderTheaterFinding(out, finding)
	}
	fmt.Fprintln(out)
}
