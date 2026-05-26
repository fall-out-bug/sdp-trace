package main

import "fmt"

func appendActionFinding(findings []string, f string, lineNo int, uses string) []string {
	if uses == "" || localUse(uses) || pinnedUse(uses) {
		return findings
	}
	return append(findings, fmt.Sprintf("workflow action is not SHA-pinned: %s:%d %s", f, lineNo, uses))
}

func unreadableActionFinding(f string) string {
	return fmt.Sprintf("workflow action pin check unreadable: %s", f)
}
