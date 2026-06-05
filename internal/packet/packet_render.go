package packet

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"
)

func RenderMarkdown(bundle Bundle) (string, error) {
	// RenderMarkdown keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	packet, err := renderablePacket(bundle)
	if err != nil {
		return "", err
	}
	return renderPacketMarkdown(packet, bundle.Manifest), nil
}

func renderablePacket(bundle Bundle) (Packet, error) {
	// renderablePacket keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	validation := Validate(bundle, time.Now().UTC())
	if validation.State != StatePass {

		return Packet{}, errors.New(strings.Join(validation.Errors, "; "))
	}
	return bundle.Packet, nil
}

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
