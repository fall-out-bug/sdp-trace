package harnessobs

import (
	"fmt"
)

func LoadValidation(path string) (Validation, error) {
	// LoadValidation keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var validation Validation
	if err := readExistingJSON(path, &validation); err != nil {
		return Validation{}, err
	}

	if validation.SchemaVersion != ValidationSchemaVersion {
		return Validation{}, fmt.Errorf("unsupported validation schema_version: %s", validation.SchemaVersion)
	}
	return validation, nil
}
