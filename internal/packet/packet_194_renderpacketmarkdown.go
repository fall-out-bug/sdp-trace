package packet

import (
	"bytes"

	"fmt"
)

func renderPacketMarkdown(packet Packet, manifest BundleManifest) string {
	// Rendering is a projection of already-validated packet and manifest data;
	// it does not reclassify evidence or compute approval.
	// Section order mirrors reviewer flow from scope to evidence to gaps.
	var out bytes.Buffer
	fmt.Fprintf(&out, "# Change Evidence Packet v0\n\n")
	fmt.Fprintf(&out, "This packet is evidence organization, not merge, release, compliance, production trust, or quality approval.\n\n")
	renderExecutiveSummary(&out, packet)
	renderMetadata(&out, packet)
	renderRows(&out, packet.Rows)
	renderTheater(&out, packet)
	renderDecisions(&out, packet.DecisionOwners)
	renderEvidence(&out, manifest)
	renderResidualGaps(&out, packet.ResidualGaps)
	renderNonProof(&out, packet)
	return out.String()
}
