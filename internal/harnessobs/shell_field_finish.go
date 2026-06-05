package harnessobs

func (scanner *shellFieldScanner) finish() []string {
	if scanner.escaped {
		scanner.current.WriteRune('\\')
	}
	scanner.flush()
	return scanner.fields
}

func (scanner *shellFieldScanner) flush() {
	if scanner.current.Len() == 0 {
		return
	}

	scanner.fields = append(scanner.fields, scanner.current.String())
	scanner.current.Reset()
}
