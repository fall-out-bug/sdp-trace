package posture

import (
	"fmt"
)

func explainMovementLines(rows []MovementRow) []string {
	return formattedLines(rows, explainMovementLine)
}

func explainMovementLine(row MovementRow) string {
	return fmt.Sprintf("movement %s %s current=%d previous=%d delta=%d comparable=%t reason=%s", row.ID, row.MetricID, row.CurrentValue, row.PreviousValue, row.Delta, row.Comparable, row.NonComparableReason)
}
