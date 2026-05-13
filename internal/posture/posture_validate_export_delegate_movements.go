package posture

func validateExportMovements(result ExportResult) error {
	return validateMovementRows(result.MovementRows)
}

func validateExportMovementSummary(result ExportResult) error {
	return validateMovementSummary(result.MovementSummary)
}

func validateExportRefusals(result ExportResult) error {
	return validateRefusalRows(result.RefusalRows)
}
