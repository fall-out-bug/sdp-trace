package managed

func managedConditions(input Input) []Condition {
	// Conditions are grouped by authority, observation, and closure so missing
	// evidence stays separate from failed evidence.

	conditions := managedAuthorityConditions(input)
	conditions = append(conditions, managedObservationConditions(input)...)
	conditions = append(conditions, managedClosureConditions(input)...)
	return conditions
}

func managedAuthorityConditions(input Input) []Condition {
	// Authority conditions decide whether the selected adapter was permitted before
	// observation evidence can be trusted.

	return []Condition{
		pass("managed_profile_explicitly_selected", "managed_profile_selected", "managed harness profile was explicitly selected"),
		policyCondition(input.Policy),
		registryCondition(input.Registry),
		boundaryCondition(input),
		adapterIdentityCondition(input),
		capabilityCondition(input),
	}
}

func managedObservationConditions(input Input) []Condition {
	// Observation conditions check adapter activity, event coverage, suppression,
	// bypasses, and witness data without upgrading missing observations.

	return []Condition{
		adapterActivationCondition(input),
		adapterConnectionCondition(input),
		eventGroupCondition(input, "required_harness_events_observed", "harness"),
		eventGroupCondition(input, "required_tool_events_observed", "tool"),
		eventGroupCondition(input, "required_file_mutations_observed", "file"),
		testProvenanceCondition(input),
		suppressionCondition(input),
		bypassCondition(input),
	}
}

func managedClosureConditions(input Input) []Condition {
	// Closure conditions bind override and witness evidence after observation checks
	// have identified the adapter and event surface.

	return []Condition{
		witnessCondition(input),
		overrideCondition(input),
	}
}

func managedAssessmentResult(conditions []Condition) AssessmentResult {
	// Result assembly keeps condition states, reasons, and next actions tied to the
	// machine evidence emitted by each managed gate.

	return AssessmentResult{
		SchemaVersion:            SchemaVersion,
		SelectedProfile:          ProfileManagedHarness,
		ManagedHarnessAssessment: topLevel(conditions),
		TrustScope:               TrustScopeManaged,
		ManagedConditions:        conditions,
	}
}
