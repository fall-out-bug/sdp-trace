package posture

func ValidateExportResult(result ExportResult) error {
	// ValidateExportResult keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	for _, validate := range exportResultValidators {
		if err := validate(result); err != nil {
			return err
		}
	}
	return nil
}

var exportResultValidators = []func(ExportResult) error{
	validateExportHeader,
	validateExportCollections,
	validateExportInputSelections,
	validateExportMetrics,
	validateExportMovements,
	validateExportMovementSummary,
	validateExportRefusals,
}
