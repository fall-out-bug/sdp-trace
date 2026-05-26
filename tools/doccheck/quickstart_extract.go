package main

import "strings"

func quickstartCommands(doc string) []string {
	var commands []string
	inCodeBlock := false
	for _, line := range strings.Split(doc, "\n") {
		var cmd string
		inCodeBlock, cmd = processQuickstartLine(inCodeBlock, line)
		if cmd != "" {
			commands = append(commands, cmd)
		}
	}
	return uniqueSorted(commands)
}
