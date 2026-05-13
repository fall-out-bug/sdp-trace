package posture

func validateExportInputSelections(result ExportResult) error {
	return validateInputSelectionRows(result.InputSelection)
}

func validateExportMetrics(result ExportResult) error {
	return validateMetricRows(result.MetricRows)
}
