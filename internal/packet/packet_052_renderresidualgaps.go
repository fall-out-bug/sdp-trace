package packet

import (
	"bytes"

	"fmt"
)

func renderResidualGaps(out *bytes.Buffer, gaps []ResidualGap) {
	// renderResidualGaps keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	fmt.Fprintf(out, "## Residual Gaps\n\n")
	if len(gaps) == 0 {
		renderNoResidualGaps(out)
		return
	}
	renderResidualGapRows(out, gaps)
}
