package posture

import (
	"errors"
	"time"
)

func validateExportHeader(result ExportResult) error {
	// validateExportHeader keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	if err := validateExportSchemaHeader(result); err != nil {
		return err
	}
	return validateExportGeneratedAt(result.GeneratedAt)
}

func validateExportSchemaHeader(result ExportResult) error {
	// validateExportSchemaHeader keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	if unsupportedExportHeader(result) {
		return errors.New("unsupported posture export")
	}
	if malformedExportHeader(result) {
		return errors.New("malformed posture export")
	}
	return nil
}

func validateExportGeneratedAt(generatedAt string) error {
	// validateExportGeneratedAt keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	if _, err := time.Parse(time.RFC3339, generatedAt); err != nil {
		return errors.New("malformed posture export generated_at")
	}
	return nil
}
