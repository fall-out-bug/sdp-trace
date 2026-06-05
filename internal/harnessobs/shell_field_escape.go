package harnessobs

func (scanner *shellFieldScanner) consumeEscaped(r rune) bool {
	if !scanner.escaped {
		return false
	}
	if r != '\n' {
		scanner.current.WriteRune('\\')
		scanner.current.WriteRune(r)
	}
	scanner.escaped = false
	return true
}

func (scanner *shellFieldScanner) startsEscape(r rune) bool {
	return scanner.quote != '\'' && r == '\\'
}
