package packet

import "strings"

func appendPromptBoundaryFinding(packet Packet, boundary PromptBoundary) Packet {
	// appendPromptBoundaryFinding keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	classification := ClassifyPromptBoundary(boundary)
	if classification.Verdict == "contaminated" {

		packet.TheaterFindings = append(packet.TheaterFindings, TheaterFinding{
			ReasonCode:          "prompt_contamination",
			State:               StateFail,
			Severity:            "P0",
			Finding:             strings.Join(classification.Reasons, "; "),
			TriggerEvidenceRefs: []string{"prompt:boundary"},
		})
	}
	return packet
}
