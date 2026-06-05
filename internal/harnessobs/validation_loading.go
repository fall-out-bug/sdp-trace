package harnessobs

import "fmt"

// Validation loading reads a rendered validation artifact and enforces the
// schema version before callers trust its shape.
func LoadValidation(path string) (Validation, error) {
	var validation Validation
	if err := readExistingJSON(path, &validation); err != nil {
		return Validation{}, err
	}
	if validation.SchemaVersion != ValidationSchemaVersion {
		return Validation{}, fmt.Errorf("unsupported validation schema_version: %s", validation.SchemaVersion)
	}
	return validation, nil
}
