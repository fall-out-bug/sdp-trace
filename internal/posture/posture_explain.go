package posture

import (
	"fmt"
	"strings"
)

func Explain(result ExportResult) (string, error) {
	// Explanation output is derived from structured rows only, then checked
	// again so renderer text cannot leak unsafe payload classes.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	var lines []string
	lines = append(lines, explainHeaderLines(result)...)
	lines = append(lines, fmt.Sprintf("movement_summary comparable=%d non_comparable=%d", result.MovementSummary.ComparableCount, result.MovementSummary.NonComparableCount))
	lines = append(lines, explainRefusalLines(result.RefusalRows)...)
	lines = append(lines, explainMetricLines(result.MetricRows)...)
	lines = append(lines, explainMovementLines(result.MovementRows)...)
	lines = append(lines, explainOutputSafetyLines(result.OutputSafety.VerifiedAbsentSensitiveClasses)...)

	rendered := strings.Join(lines, "\n") + "\n"
	if unsafeOutput(rendered) {
		return "", fmt.Errorf("output_safety_violation")
	}
	return rendered, nil
}
