package main

import "strings"

func documentedCommands(doc string) []string {
	// The current command surface section is intentionally simple markdown so a
	// small parser can keep it honest in CI without Node or shell-specific code.
	section := currentCommandSurface(doc)
	commands := make([]string, 0, strings.Count(section, "- `sdp-trace "))
	for _, line := range strings.Split(section, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "- `sdp-trace ") || !strings.HasSuffix(line, "`") {
			continue
		}
		commands = append(commands, strings.TrimSuffix(strings.TrimPrefix(line, "- `"), "`"))
	}
	return uniqueSorted(commands)
}

func currentCommandSurface(doc string) string {
	// Bound the comparison to the canonical list; later examples may contain
	// narrative snippets that are not meant to be exhaustive.
	_, after, ok := strings.Cut(doc, "Current command surface:")
	if !ok {
		return ""
	}
	before, _, ok := strings.Cut(after, "Do not add aliases")
	if !ok {
		return after
	}
	return before
}
