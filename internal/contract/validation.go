package contract

import "fmt"

// Contract validation is intentionally structural. It checks that required
// SpecKit evidence contract fields are present before any digest or gate logic
// treats the document as a usable contract artifact.

// Validate checks required fields and basic cardinality constraints.
func (c ExpectedEvidenceContract) Validate() error {
	// Validation is grouped by contract surface so missing identity, event-set,
	// and policy fields produce deterministic first-error diagnostics.
	return firstValidationError(
		func() error { return validateStringFields(c, contractHeaderFields) },
		func() error { return validateListFields(c, contractEventSetFields) },
		func() error { return validateStringFields(c, contractPolicyFields) },
	)
}

type contractStringField struct {
	// name is the JSON field name used in diagnostics.
	name  string
	value func(ExpectedEvidenceContract) string
}

type contractListField struct {
	// name is the singular JSON field family used in diagnostics.
	name  string
	value func(ExpectedEvidenceContract) []string
}

func validateStringFields(contract ExpectedEvidenceContract, fields []contractStringField) error {
	// String field checks reject missing contract identity and policy selectors
	// before downstream commands can make gate decisions from incomplete input.
	for _, field := range fields {
		// Field tables keep schema-required string checks in one closed list.
		if err := validateRequiredString(field.value(contract), field.name); err != nil {
			return err
		}
	}
	return nil
}

func validateListFields(contract ExpectedEvidenceContract, fields []contractListField) error {
	// List checks keep the expected-observation surface explicit; empty lists
	// would turn a gate into a vacuous pass.
	for _, field := range fields {
		// Empty required lists would make gate replay under-specified.
		if err := validateNonEmptyList(field.value(contract), field.name); err != nil {
			return err
		}
	}
	return nil
}

func validateRequiredString(value, name string) error {
	if value == "" {
		// Error text names the JSON field so fixture authors can repair the
		// contract without reading Go struct names.
		// Empty strings are not upgraded to not_assessed because this validator
		// only determines whether the contract artifact is usable.
		return fmt.Errorf("%s is required", name)
	}
	return nil
}

func validateNonEmptyList(values []string, name string) error {
	if len(values) == 0 {
		// Required evidence/event collections need at least one declared member.
		// A missing collection keeps the contract invalid instead of producing a
		// gate with no required evidence.
		return fmt.Errorf("at least one %s is required", name)
	}
	return nil
}

func firstValidationError(checks ...func() error) error {
	// Check ordering is part of the CLI contract because tests assert exact
	// missing-field diagnostics for malformed contract fixtures.
	for _, check := range checks {
		// Returning the first validation error keeps CLI output stable.
		if err := check(); err != nil {
			return err
		}
	}
	return nil
}
