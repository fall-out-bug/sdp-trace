package harnessobs

func (scanner *shellFieldScanner) consumeQuoted(r rune) bool {
	if scanner.quote == 0 {
		return false
	}
	if r == scanner.quote {
		scanner.quote = 0
		return true
	}
	scanner.current.WriteRune(r)
	return true
}

func shellQuote(r rune) bool {
	return r == '\'' || r == '"'
}
