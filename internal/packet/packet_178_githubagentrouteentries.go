package packet

import (
	"strings"
)

func githubAgentRouteEntries(input GitHubPREvidenceInput) []BundleEntry {
	// githubAgentRouteEntries keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if len(input.AgentRouteRefs) > 0 {
		entry := bundleEntry("agent:route", "harness", strings.Join(input.AgentRouteRefs, ", "), "external_ref")
		if strings.TrimSpace(input.AgentRouteDigest) != "" {

			entry.Digest = input.AgentRouteDigest
		}
		entry.EvidenceKind = input.AgentRouteEvidenceKind
		entry.ObservedComponents = input.AgentRouteComponents
		entry = authorityEntry(entry, "recorder", "recorder_owned", "sdp-trace recorder run", "external_retained_artifact", input.AgentRouteDigest)
		return []BundleEntry{entry}
	}
	return nil
}
