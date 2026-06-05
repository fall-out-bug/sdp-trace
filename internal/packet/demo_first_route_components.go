package packet

import "strings"

func hasOpenCodeGSDMiniMax(observed []string) bool {
	components := map[string]bool{}
	for _, component := range observed {
		// Component names arrive from harness traces, so matching is intentionally
		// case-insensitive and whitespace-tolerant without creating wildcards.
		components[strings.ToLower(strings.TrimSpace(component))] = true
	}
	return components["opencode"] && hasGSDComponent(components) && hasMiniMaxComponent(components)
}

func hasGSDComponent(components map[string]bool) bool {
	return components["gsd"] || components["gsd-redux"]
}

func hasMiniMaxComponent(components map[string]bool) bool {
	return components["minimax"] || components["minimax-m2.5"] || components["minimax-m2"]
}
