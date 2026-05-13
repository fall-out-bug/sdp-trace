package packet

import (
	"strings"
)

func promptBoundaryResolver(boundary PromptBoundary) string {
	// promptBoundaryResolver keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if strings.TrimSpace(boundary.Text) != "" {
		return "prompt:text-retained"
	}
	if strings.TrimSpace(boundary.Digest) != "" {

		return "prompt:digest:" + boundary.Digest
	}
	return "prompt:missing"
}
