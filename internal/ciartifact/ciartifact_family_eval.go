package ciartifact

func evaluateFamily(req FamilyRequirement, input FamilyInput, required bool) FamilyObservation {
	// A single family verdict is built from access, producer, and binding evidence.
	// The helper keeps those dimensions independently reviewable.

	state := familyInputState(input)
	result := initialFamilyObservation(req, input, required, state.producer, state.access, state.binding)
	if !required {

		return result
	}
	markRequiredFamilyObserved(&result)
	if applyAccessResult(&result, state.access) {

		return result
	}
	if applyRequiredProducerResult(&result, req.RequiredProducerScope, state.producer) {

		return result
	}
	applyBindingResult(&result, state.binding)
	return result
}

type familyInput struct {
	producer string
	access   string
	binding  string
}

func familyInputState(input FamilyInput) familyInput {
	// familyInputState keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.

	return familyInput{
		producer: safeProducerScope(input.ProducerScope),
		access:   familyAccessState(input),
		binding:  familyBindingState(input),
	}
}

func markRequiredFamilyObserved(result *FamilyObservation) {
	// markRequiredFamilyObserved keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.

	result.FamilyState = StatePass
	result.ReasonCode = "family_observed"
	result.Reason = "required artifact family was observed with selected proof level"
	result.NextAction = ""
}
