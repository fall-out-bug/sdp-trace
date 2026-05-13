package authority

var aggregateStatePriority = map[string]int{
	StateCannotVerify:     3,
	StateOutsideAuthority: 2,
	StateWithinAuthority:  1,
	StateNotAssessed:      0,
}

var aggregateStateByRank = []string{
	StateNotAssessed,
	StateWithinAuthority,
	StateOutsideAuthority,
	StateCannotVerify,
}
