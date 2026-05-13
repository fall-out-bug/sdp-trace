package packet

func theaterFindingState(state string) bool {
	return state == StatePartial || state == StateFail || state == StateCannotVerify
}
