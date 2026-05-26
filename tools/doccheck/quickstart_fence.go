package main

import "strings"

func openingFence(line string) (bool, string) {
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "```") && len(trimmed) > 3 {
		return true, ""
	}
	return false, ""
}

func closingFence(line string) bool {
	return strings.TrimSpace(line) == "```"
}

// isQuickstartCommand requires a trailing space after "sdp-trace" so that a
// bare "go run ./cmd/sdp-trace" line (without a subcommand) is ignored. This
// is intentional: --help is handled as an exact-match meta-flag elsewhere.
func isQuickstartCommand(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "go run ./cmd/sdp-trace ")
}
