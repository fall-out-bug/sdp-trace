package packet

func passOrPartial(state string) bool {
	return state == StatePass || state == StatePartial
}
