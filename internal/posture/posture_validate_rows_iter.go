package posture

func validateInputSelectionRows(rows []InputSelection) error {
	return validateRows(rows, validateInputSelectionRow)
}

func validateMetricRows(rows []MetricRow) error {
	return validateRows(rows, validateMetricRow)
}

func validateMovementRows(rows []MovementRow) error {
	return validateRows(rows, validateMovementRow)
}

func validateRefusalRows(rows []RefusalRow) error {
	return validateRows(rows, validateRefusalRow)
}

func validateRows[T any](rows []T, validate func(T) error) error {
	// validateRows keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	for _, row := range rows {
		if err := validate(row); err != nil {
			return err
		}
	}
	return nil
}
