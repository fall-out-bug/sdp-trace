package harnessobs

var stateRank = map[string]int{
	StateFail:         4,
	StateCannotVerify: 3,
	StateNotAssessed:  2,
	StatePass:         1,
}
