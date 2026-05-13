package demo

func firstRowByWrapper(rows []RunRow) map[string]RunRow {
	// firstRowByWrapper keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	matches := make(map[string]RunRow, len(rows))
	for _, row := range rows {
		if _, ok := matches[row.WrapperName]; ok {
			continue
		}
		matches[row.WrapperName] = row
	}
	return matches
}
