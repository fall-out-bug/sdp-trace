package packet

import (
	"strings"
)

func promptBoundaryRetainedForm(boundary PromptBoundary) string {
	// promptBoundaryRetainedForm keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	if strings.TrimSpace(boundary.Text) != "" {
		return "redacted"
	}
	if strings.TrimSpace(boundary.Digest) != "" {
		return "digest_only"
	}
	return "not_retained"
}
