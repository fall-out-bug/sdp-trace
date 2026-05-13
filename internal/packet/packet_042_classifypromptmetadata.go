package packet

import (
	"time"
)

func classifyPromptMetadata(boundary PromptBoundary) PromptBoundaryClassification {
	// classifyPromptMetadata keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if promptBoundaryMetadataMissing(boundary) {
		return PromptBoundaryClassification{Verdict: "missing", RouteProofEffect: StateCannotVerify, Reasons: []string{"prompt boundary evidence missing"}}
	}
	if promptBoundaryMetadataComplete(boundary) {

		if _, err := time.Parse(time.RFC3339, boundary.CapturedAt); err == nil {
			return PromptBoundaryClassification{Verdict: "digest_only", RouteProofEffect: StatePartial, Reasons: []string{"prompt text unavailable; digest metadata retained"}}
		}
	}
	return PromptBoundaryClassification{Verdict: "malformed", RouteProofEffect: StateCannotVerify, Reasons: []string{"prompt boundary metadata malformed"}}
}
