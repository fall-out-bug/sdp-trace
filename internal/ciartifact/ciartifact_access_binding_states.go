package ciartifact

func safeAccessState(value string) string {
	// safeAccessState keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if validAccessStates[value] {
		return value
	}

	return AccessCannotVerify
}

var validAccessStates = map[string]bool{
	AccessPresent:      true,
	AccessAbsent:       true,
	AccessPartial:      true,
	AccessExpired:      true,
	AccessInaccessible: true,
	AccessMalformed:    true,
	AccessUnsafe:       true,
	AccessNotAssessed:  true,
	AccessCannotVerify: true,
}

func safeBindingState(value string) string {
	// safeBindingState keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if validBindingStates[value] {
		return value
	}

	return BindingUnverifiable
}

var validBindingStates = map[string]bool{
	BindingMatched:      true,
	BindingMismatch:     true,
	BindingAbsent:       true,
	BindingUnverifiable: true,
	BindingNotAssessed:  true,
}
