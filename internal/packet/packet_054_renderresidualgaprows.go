package packet

import (
	"bytes"

	"fmt"
)

func renderResidualGapRows(out *bytes.Buffer, gaps []ResidualGap) {
	// renderResidualGapRows keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	fmt.Fprintf(out, "| row id | state | reason | closure evidence |\n| --- | --- | --- | --- |\n")
	for _, gap := range gaps {
		fmt.Fprintf(out, "| %s | %s | %s | %s |\n", gap.RowID, gap.State, md(gap.Reason), md(gap.ClosureEvidence))
	}
	fmt.Fprintln(out)
}
