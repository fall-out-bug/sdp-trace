package ciartifact

func familyAccessState(input FamilyInput) string {
	// familyAccessState keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if input.ArtifactAccessState == "" {

		return AccessAbsent
	}
	return safeAccessState(input.ArtifactAccessState)
}

func familyBindingState(input FamilyInput) string {
	// familyBindingState keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if input.BindingState == "" {

		return BindingAbsent
	}
	return safeBindingState(input.BindingState)
}

func initialFamilyObservation(req FamilyRequirement, input FamilyInput, required bool, producer, access, binding string) FamilyObservation {
	// initialFamilyObservation keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	family := req.Family
	if family == "" {

		family = input.Family
	}

	return FamilyObservation{
		Family:              family,
		Required:            required,
		RequiredProducer:    req.RequiredProducerScope,
		ProducerScope:       producer,
		ArtifactAccessState: access,
		BindingState:        binding,
		FamilyState:         StateNotAssessed,
		ReasonCode:          "family_not_selected",
		Reason:              "artifact family was outside the selected profile scope",
	}
}
