package main

import (
	"fmt"
	"strings"
)

func compareCommandSurface(help, doc string) error {
	// Compare both directions: missing rows hide new commands, while stale rows
	// document commands that the live CLI no longer accepts.
	helpCommands := helpCommands(help)
	docCommands := documentedCommands(doc)
	missing := missingCommands(helpCommands, docCommands)
	stale := missingCommands(docCommands, helpCommands)
	if len(missing) == 0 && len(stale) == 0 {
		return nil
	}
	return fmt.Errorf("agent entrypoint command surface drift: %s", commandSurfaceDrift(missing, stale))
}

func commandSurfaceDrift(missing, stale []string) string {
	// Missing and stale rows are rendered separately so command additions and
	// removals can be fixed independently from one CI failure.
	var parts []string
	if len(missing) > 0 {
		parts = append(parts, "missing documented commands: "+strings.Join(missing, "; "))
	}
	if len(stale) > 0 {
		parts = append(parts, "stale documented commands: "+strings.Join(stale, "; "))
	}
	return strings.Join(parts, " | ")
}
