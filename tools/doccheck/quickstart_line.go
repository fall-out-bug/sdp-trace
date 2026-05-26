package main

import "strings"

func processQuickstartLine(inCodeBlock bool, line string) (bool, string) {
	if !inCodeBlock {
		return openingFence(line)
	}
	if closingFence(line) {
		return false, ""
	}
	if isQuickstartCommand(line) {
		return inCodeBlock, strings.TrimSpace(line)
	}
	return inCodeBlock, ""
}
