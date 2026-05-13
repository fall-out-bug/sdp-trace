package harnessobs

func (scanner *shellFieldScanner) startsEscape(r rune) bool {
	return scanner.quote != '\'' && r == '\\'
}
