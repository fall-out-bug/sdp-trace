package packet

import (
	"fmt"
	"time"
)

func BuildFromGitHubInput(input GitHubPREvidenceInput, generatedAt time.Time) Bundle {
	// BuildFromGitHubInput keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	packetID := fmt.Sprintf("github-pr-%d-change-evidence-packet", input.PR.Number)
	bundleID := fmt.Sprintf("%s-bundle", packetID)

	entries := githubEntries(input)
	rows := githubRows(input)
	packet := githubPacket(input, generatedAt, packetID, bundleID, rows)
	packet = appendPromptBoundaryFinding(packet, input.PromptBoundary)
	if len(input.IntegrationActions) > 0 {

		if packet.Extensions == nil {
			packet.Extensions = map[string]any{}
		}
		packet.Extensions["integration_actions"] = input.IntegrationActions
	}
	return Bundle{Packet: packet, Manifest: githubBundleManifest(bundleID, packet, entries)}
}
