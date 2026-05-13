package posture

func validateMovementRow(row MovementRow) error {
	return malformedRowError(malformedMovementRow(row), "malformed posture export movement_row")
}

func malformedMovementRow(row MovementRow) bool {
	return malformedMovementIdentity(row) || malformedMovementValues(row) || malformedMovementComparison(row)
}
