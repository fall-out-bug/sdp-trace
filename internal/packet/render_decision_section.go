package packet

import (
	"bytes"
	"fmt"
)

func renderDecisions(out *bytes.Buffer, owners []DecisionOwner) {
	// renderDecisions keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	fmt.Fprintf(out, "## Decision Ownership\n\n")
	fmt.Fprintf(out, "| decision | owner | state | reason |\n| --- | --- | --- | --- |\n")
	for _, owner := range owners {
		fmt.Fprintf(out, "| %s | %s | %s | %s |\n", md(owner.Decision), md(owner.Owner), owner.State, md(owner.Reason))
	}
	fmt.Fprintln(out)
}
