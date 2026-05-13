package packet

import (
	"bytes"

	"fmt"
	"sort"
	"strings"
)

func renderRows(out *bytes.Buffer, rows []Row) {
	// renderRows keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	fmt.Fprintf(out, "## Required Rows\n\n")
	fmt.Fprintf(out, "| row id | state | answer | evidence refs | gap / next evidence | owner |\n| --- | --- | --- | --- | --- | --- |\n")
	ordered := append([]Row(nil), rows...)

	sort.SliceStable(ordered, func(i, j int) bool {
		return requiredRowIndex(ordered[i].ID) < requiredRowIndex(ordered[j].ID)
	})
	for _, row := range ordered {
		gap := row.Reason
		if gap == "" {
			gap = "none"
		}
		fmt.Fprintf(out, "| %s | %s | %s | %s | %s | %s |\n", row.ID, row.State, md(row.Summary), md(strings.Join(row.EvidenceRefs, ", ")), md(gap), md(row.Owner))
	}
	fmt.Fprintln(out)
}
