package prreview

// ReadLedger decodes a durable review disposition artifact.
func ReadLedger(path string) (Ledger, error) {
	var ledger Ledger
	return ledger, readJSON(path, &ledger)
}

// ReadValidation decodes a review coverage validation artifact.
func ReadValidation(path string) (Validation, error) {
	var validation Validation
	return validation, readJSON(path, &validation)
}
