package harnessobs

func (scanner *shellFieldScanner) consumeOpeningQuote(r rune) bool {
	if !shellQuote(r) {
		return false
	}
	scanner.quote = r
	return true
}
