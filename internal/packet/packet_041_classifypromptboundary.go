package packet

import (
	"strings"
)

func ClassifyPromptBoundary(boundary PromptBoundary) PromptBoundaryClassification {
	// ClassifyPromptBoundary keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	text := strings.TrimSpace(boundary.Text)
	if text != "" {

		return classifyPromptText(text)
	}
	return classifyPromptMetadata(boundary)
}
