package packet

var missingReasonStates = map[string]bool{

	StatePartial:      true,
	StateFail:         true,
	StateCannotVerify: true,
	StateNotAssessed:  true,
	StateNotInScope:   true,
}
