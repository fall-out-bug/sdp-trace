package posture

func formattedLines[T any](rows []T, format func(T) string) []string {
	// formattedLines keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	lines := make([]string, 0, len(rows))
	for _, row := range rows {
		lines = append(lines, format(row))
	}
	return lines
}
