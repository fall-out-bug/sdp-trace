package posture

import (
	"fmt"
)

func explainMetricLines(rows []MetricRow) []string {
	return formattedLines(rows, explainMetricLine)
}

func explainMetricLine(row MetricRow) string {
	return fmt.Sprintf("metric %s %s numerator=%d denominator=%d window=%s dimension_key=%s not_assessed_count=%d", row.ID, row.MetricID, row.Numerator, row.Denominator, row.TimeWindow, row.DimensionKey, row.NotAssessedCount)
}
