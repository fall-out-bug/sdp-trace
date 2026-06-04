package packet

import "time"

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
