package forensic

func rawReferenceCondition(input Input) Condition {
	// rawReferenceCondition keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	for _, event := range input.Run.Events {
		ref := event.RawReference
		if ref == nil {
			continue
		}

		if condition, ok := validateRawReference(ref); !ok {
			return condition
		}
	}
	return pass("raw_reference_bound", "raw_reference_bound", "raw references are digest-bound and access-verifiable")
}

func validateRawReference(ref *RawReference) (Condition, bool) {
	condition, ok := rawReferenceValidationFailure(ref)
	return condition, ok
}

func rawReferenceValidationFailure(ref *RawReference) (Condition, bool) {
	// rawReferenceValidationFailure keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	for _, rule := range rawReferenceValidationRules {
		if rule.invalid(ref) {

			return rule.condition, false
		}
	}
	return Condition{}, true
}
