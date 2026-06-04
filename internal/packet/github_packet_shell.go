package packet

import (
	"fmt"
	"time"
)

// Source-change projection binds the generated packet to the GitHub PR and the
// explicit commit range. It does not infer repository identity outside the
// retained input.
func githubSourceChange(input GitHubPREvidenceInput) SourceChange {
	return SourceChange{
		Repository:  input.PR.URL,
		ChangeID:    fmt.Sprintf("PR-%d", input.PR.Number),
		URL:         input.PR.URL,
		BaseRef:     input.PR.BaseRef,
		HeadRef:     input.PR.HeadRef,
		CommitRange: input.CommitRange.Base + ".." + input.CommitRange.Head,
		HeadSHA:     input.PR.HeadSHA,
	}
}

// GitHub packet projection assembles generated rows with packet metadata and
// leaves residual gaps, decision owners, and digest material to their dedicated
// helpers.
func githubPacket(input GitHubPREvidenceInput, generatedAt time.Time, packetID, bundleID string, rows []Row) Packet {
	// githubPacket keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	return Packet{
		PacketVersion:   PacketSchemaVersion,
		PacketID:        packetID,
		SourceChange:    githubSourceChange(input),
		GeneratedAt:     generatedAt.UTC().Format(time.RFC3339),
		AuthoringMethod: AuthoringToolGenerated,
		SelectedProfile: "change-host-rich-v0",
		RedactionPolicy: "not_assessed",
		BundleRef:       bundleID,
		PacketState:     "draft",
		Projection:      Projection{Kind: ProjectionCanonical, Canonical: true, ArtifactRef: "packet:markdown"},
		Rows:            rows,
		ResidualGaps:    residualGapsForRows(rows),
		DecisionOwners:  defaultDecisionOwners(),
		NonApproval:     "This packet does not approve merge, release, compliance, production trust, semantic correctness, or signed external trust.",
	}
}
