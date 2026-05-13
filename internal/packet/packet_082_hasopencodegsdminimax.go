package packet

import (
	"strings"
)

func hasOpenCodeGSDMiniMax(observed []string) bool {
	// hasOpenCodeGSDMiniMax keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	components := map[string]bool{}
	for _, component := range observed {

		components[strings.ToLower(strings.TrimSpace(component))] = true
	}
	return components["opencode"] && components["gsd"] && hasMiniMaxComponent(components)
}
