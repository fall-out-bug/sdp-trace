package harnessobs

func shellQuote(r rune) bool {
	return r == '\'' || r == '"'
}
