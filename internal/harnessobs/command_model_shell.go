package harnessobs

import "path/filepath"

func shellCommandString(command []string) string {
	if !shellCommandShape(command) {
		return ""
	}
	base := filepath.Base(command[0])
	if base != "sh" && base != "bash" {
		return ""
	}

	return command[2]
}

func shellCommandShape(command []string) bool {
	return len(command) >= 3 && command[1] == "-c"
}
