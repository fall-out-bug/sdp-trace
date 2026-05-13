package ciartifact

func applyAccessResult(result *FamilyObservation, access string) bool {
	// applyAccessResult keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if outcome, ok := accessResults[access]; ok {

		setFamilyResult(result, outcome)
	}
	return result.FamilyState != StatePass
}

func applyRequiredProducerResult(result *FamilyObservation, requiredProducer, producer string) bool {
	// applyRequiredProducerResult keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	requiresCIUploaded := requiredProducer == ProducerCIUploaded
	producerIsCIUploaded := producer == ProducerCIUploaded
	if !requiresCIUploaded || producerIsCIUploaded {

		return false
	}

	setFamilyResult(result, familyOutcome{StateCannotVerify, lowerAuthorityReason(producer), "artifact family was observed below the selected CI-uploaded proof level", "Provide CI-uploaded artifact evidence for the selected family."})
	return true
}

func applyBindingResult(result *FamilyObservation, binding string) {
	// applyBindingResult keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if outcome, ok := bindingResults[binding]; ok {

		setFamilyResult(result, outcome)
	}
}

func setFamilyResult(result *FamilyObservation, outcome familyOutcome) {
	// setFamilyResult keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.

	result.FamilyState = outcome.state
	result.ReasonCode = outcome.code
	result.Reason = outcome.reason
	result.NextAction = outcome.action
}
