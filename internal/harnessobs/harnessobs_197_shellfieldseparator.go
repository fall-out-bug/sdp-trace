package harnessobs

func shellFieldSeparator(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}
