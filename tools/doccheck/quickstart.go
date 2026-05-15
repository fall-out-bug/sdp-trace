package main

import (
	"fmt"
	"strings"
)

const contributorQuickstart = "docs/contributor-quickstart.md"

// requiredQuickstartCommands are the minimal commands that must appear in the
// contributor quick start so a new reader can verify the local environment.
var requiredQuickstartCommands = []string{
	"go run ./cmd/sdp-trace --help",
	"go run ./cmd/sdp-trace doctor",
	"go run ./cmd/sdp-trace wrap",
	"go run ./cmd/sdp-trace verify",
	"go run ./cmd/sdp-trace explain",
}

func compareQuickstartWithRegistry(quickstart string) error {
	registry, err := registryUsages()
	if err != nil {
		return err
	}
	qsCmds := quickstartCommands(quickstart)
	missing := missingQuickstartCommands(qsCmds)
	stale := staleQuickstartCommands(qsCmds, registry)
	if len(missing) == 0 && len(stale) == 0 {
		return nil
	}
	return fmt.Errorf("quickstart drift: %s", quickstartDrift(missing, stale))
}

func missingQuickstartCommands(qsCmds []string) []string {
	qsSet := stringSliceToSet(qsCmds)
	var missing []string
	for _, req := range requiredQuickstartCommands {
		if !quickstartHasCommand(qsSet, req) {
			missing = append(missing, req)
		}
	}
	return missing
}

func quickstartHasCommand(qsSet map[string]bool, req string) bool {
	if req == "go run ./cmd/sdp-trace --help" {
		_, ok := qsSet["go run ./cmd/sdp-trace --help"]
		return ok
	}
	return setContainsPrefix(qsSet, req)
}

func setContainsPrefix(set map[string]bool, prefix string) bool {
	for s := range set {
		if strings.Contains(s, prefix) {
			return true
		}
	}
	return false
}

func staleQuickstartCommands(qsCmds, registry []string) []string {
	registrySet := stringSliceToSet(registry)
	var stale []string
	for _, qs := range qsCmds {
		if qs == "go run ./cmd/sdp-trace --help" {
			continue // meta-flag, not a registry command
		}
		if !isKnownCommand(qs, registrySet) {
			stale = append(stale, qs)
		}
	}
	return stale
}

func isKnownCommand(qs string, registrySet map[string]bool) bool {
	normalized := normalizeQuickstartCommand(qs)
	if registrySet[normalized] {
		return true
	}
	return registryHasBase(registrySet, baseCommand(normalized))
}

func normalizeQuickstartCommand(qs string) string {
	const prefix = "go run ./cmd/sdp-trace "
	if strings.HasPrefix(qs, prefix) {
		return "sdp-trace " + strings.TrimPrefix(qs, prefix)
	}
	return qs
}

func registryHasBase(registrySet map[string]bool, base string) bool {
	for reg := range registrySet {
		if baseCommand(reg) == base {
			return true
		}
	}
	return false
}

func baseCommand(usage string) string {
	fields := strings.Fields(usage)
	if len(fields) >= 2 && fields[0] == "sdp-trace" {
		return fields[1]
	}
	return ""
}

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

func processQuickstartLine(inCodeBlock bool, line string) (bool, string) {
	if isCodeFence(line) {
		return !inCodeBlock, ""
	}
	if inCodeBlock && isQuickstartCommand(line) {
		return inCodeBlock, strings.TrimSpace(line)
	}
	return inCodeBlock, ""
}

func isCodeFence(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "```")
}

func isQuickstartCommand(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "go run ./cmd/sdp-trace ")
}

func quickstartDrift(missing, stale []string) string {
	var parts []string
	if len(missing) > 0 {
		parts = append(parts, "missing required commands: "+strings.Join(missing, "; "))
	}
	if len(stale) > 0 {
		parts = append(parts, "stale commands: "+strings.Join(stale, "; "))
	}
	return strings.Join(parts, " | ")
}
