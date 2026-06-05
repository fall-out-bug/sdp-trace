package packet

import "strings"

func gapForRowWithClosure(gaps []ResidualGap, rowID string) bool {
	for _, gap := range gaps {
		// A matching residual-gap row is not enough by itself; closure evidence
		// is the replayable next step that keeps cannot_verify honest.
		if gap.RowID == rowID && strings.TrimSpace(gap.ClosureEvidence) != "" {
			return true
		}
	}
	return false
}
