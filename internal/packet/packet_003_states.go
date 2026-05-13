package packet

var states = map[string]bool{
	StatePass:         true,
	StatePartial:      true,
	StateFail:         true,
	StateCannotVerify: true,
	StateNotAssessed:  true,
	StateNotInScope:   true,
}
