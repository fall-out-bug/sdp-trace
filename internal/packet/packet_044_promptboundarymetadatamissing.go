package packet

import (
	"strings"
)

func promptBoundaryMetadataMissing(boundary PromptBoundary) bool {
	// promptBoundaryMetadataMissing keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	return strings.TrimSpace(boundary.Digest) == "" &&
		strings.TrimSpace(boundary.CaptureActor) == "" &&
		strings.TrimSpace(boundary.CapturedAt) == "" &&
		strings.TrimSpace(boundary.CaptureMethod) == ""
}
