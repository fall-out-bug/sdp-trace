package posture

import (
	"fmt"
	"strings"
)

func movementRowForKey(index int, key string, rowsByWindow map[string]MetricRow, currentWindow, previousWindow string) MovementRow {
	// movementRowForKey keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	parts := strings.Split(key, "|")
	current, hasCurrent := rowsByWindow[currentWindow]
	previous, hasPrevious := rowsByWindow[previousWindow]

	row := MovementRow{
		ID:              fmt.Sprintf("movement.%04d", index),
		MetricID:        parts[0],
		MetricVersion:   parts[1],
		DimensionKey:    parts[2],
		ComparisonBasis: "same_profile_metric_dimension_window",
		Comparable:      hasCurrent && hasPrevious,
	}
	applyMovementWindowValues(&row, current, hasCurrent, previous, hasPrevious)
	return row
}

func applyMovementWindowValues(row *MovementRow, current MetricRow, hasCurrent bool, previous MetricRow, hasPrevious bool) {
	// applyMovementWindowValues keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	if hasCurrent {
		row.CurrentMetricRowRef = current.ID
		row.CurrentValue = current.Numerator
	}
	if hasPrevious {
		row.PreviousMetricRowRef = previous.ID
		row.PreviousValue = previous.Numerator
	}

	row.Delta = row.CurrentValue - row.PreviousValue
}
