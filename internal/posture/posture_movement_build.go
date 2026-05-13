package posture

func buildMovements(metrics []MetricRow, currentWindow, previousWindow string) ([]MovementRow, MovementSummary) {
	// buildMovements keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	byKey := metricsByMovementKey(metrics)
	keys := sortedMovementKeys(byKey)
	var rows []MovementRow
	summary := MovementSummary{NonComparableReason: map[string]int{}}
	for i, key := range keys {

		row := movementRowForKey(i+1, key, byKey[key], currentWindow, previousWindow)
		summarizeMovement(&summary, &row)
		rows = append(rows, row)
	}
	return rows, summary
}

func sortedMovementKeys(metrics map[string]map[string]MetricRow) []string {
	return sortedMapKeys(metrics)
}
