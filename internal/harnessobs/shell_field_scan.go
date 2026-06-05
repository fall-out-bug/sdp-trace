package harnessobs

func (scanner *shellFieldScanner) scan(r rune) {
	if scanner.consumeEscaped(r) {
		return
	}
	if scanner.startsEscape(r) {
		scanner.escaped = true
		return
	}
	if scanner.consumeQuoted(r) {
		return
	}

	scanner.consumeUnquoted(r)
}
