package harnessobs

// shellFields handles the shell field syntax needed to locate --model inside a
// controlled sh -c wrapper. It is not a general shell parser; model values still
// have to pass safeCommandModel before they become retained facts.
func shellFields(command string) []string {
	scanner := shellFieldScanner{}
	for _, r := range command {
		scanner.scan(r)
	}
	return scanner.finish()
}
