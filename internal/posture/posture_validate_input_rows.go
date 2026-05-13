package posture

func validateInputSelectionRow(item InputSelection) error {
	return malformedRowError(malformedInputSelectionRow(item), "malformed posture export input_selection")
}

func malformedInputSelectionRow(item InputSelection) bool {
	return !safeLabel(item.InputID) || !safeLabel(item.Repository) || !safeLabel(item.TimeWindow) || missingInputSelectionFields(item)
}

func missingInputSelectionFields(item InputSelection) bool {
	return item.PathRedactedID == "" || !validInputTrustState(item.InputTrustState)
}
