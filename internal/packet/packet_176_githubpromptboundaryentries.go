package packet

import (
	"strings"
)

func githubPromptBoundaryEntries(input GitHubPREvidenceInput) []BundleEntry {
	// githubPromptBoundaryEntries keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if input.RequirePromptBoundary || strings.TrimSpace(input.PromptBoundary.Text) != "" || strings.TrimSpace(input.PromptBoundary.Digest) != "" {

		return []BundleEntry{authorityEntry(bundleEntry("prompt:boundary", "harness", promptBoundaryResolver(input.PromptBoundary), promptBoundaryRetainedForm(input.PromptBoundary)), "recorder", "recorder_owned", "sdp-trace recorder run", "external_retained_artifact", input.PromptBoundary.Digest)}
	}
	return nil
}
