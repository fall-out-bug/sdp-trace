package packet

import "strings"

func gapForRow(gaps []ResidualGap, rowID string) bool {
	for _, gap := range gaps {
		// A residual-gap row only closes validation when it names an actual
		// reason; a blank row would hide an unresolved trust gap.
		if gap.RowID == rowID && strings.TrimSpace(gap.Reason) != "" {
			return true
		}
	}
	return false
}
