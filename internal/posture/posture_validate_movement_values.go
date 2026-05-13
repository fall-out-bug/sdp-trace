package posture

func malformedMovementValues(row MovementRow) bool {
	return row.CurrentValue < 0 || row.PreviousValue < 0
}

func malformedMovementComparison(row MovementRow) bool {
	return !validComparisonBasis(row.ComparisonBasis) || (!row.Comparable && row.NonComparableReason != "non_comparable_missing_window")
}
