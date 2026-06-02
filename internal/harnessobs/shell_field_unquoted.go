package harnessobs

import "strings"

func (scanner *shellFieldScanner) consumeUnquoted(r rune) {
	if scanner.consumeOpeningQuote(r) {
		return
	}
	if shellFieldSeparator(r) {
		scanner.flush()
		return
	}
	scanner.current.WriteRune(r)
}

func shellFieldSeparator(r rune) bool {
	return strings.ContainsRune(" \t\n\r", r)
}
