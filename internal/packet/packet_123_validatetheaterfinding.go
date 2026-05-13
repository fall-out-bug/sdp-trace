package packet

import (
	"strings"
)

func (v *bundleValidator) validateTheaterFinding(finding TheaterFinding) {
	// validateTheaterFinding keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if strings.TrimSpace(finding.ReasonCode) == "" {
		v.add("theater finding requires reason_code")
	} else if !theaterReasonCodes[finding.ReasonCode] {
		v.add("theater finding has unknown reason_code %q", finding.ReasonCode)
	}
	for _, ref := range finding.TriggerEvidenceRefs {

		v.validateEvidenceRef("theater finding "+finding.ReasonCode, StatePartial, ref)
	}
}
