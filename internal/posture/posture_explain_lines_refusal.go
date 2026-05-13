package posture

import (
	"fmt"
)

func explainRefusalLines(rows []RefusalRow) []string {
	return formattedLines(rows, explainRefusalLine)
}

func explainRefusalLine(row RefusalRow) string {
	return fmt.Sprintf("refusal %s input=%s reason=%s state=%s", row.ID, row.InputID, row.RefusalReason, row.InputTrustState)
}
