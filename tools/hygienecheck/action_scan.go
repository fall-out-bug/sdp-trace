package main

import "bufio"

func scanActionPins(scanner *bufio.Scanner, f string) []string {
	var findings []string
	for lineNo := 1; scanner.Scan(); lineNo++ {
		findings = appendActionFinding(findings, f, lineNo, usesValue(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		findings = append(findings, unreadableActionFinding(f))
	}
	return findings
}
