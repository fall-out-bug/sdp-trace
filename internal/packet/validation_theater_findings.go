package packet

import "strings"

// Theater findings are not standalone proof. Their trigger refs still need
// manifest-backed resolver evidence before a packet can explain the claim.
func (v *bundleValidator) validateTheaterFinding(finding TheaterFinding) {
	if strings.TrimSpace(finding.ReasonCode) == "" {
		v.add("theater finding requires reason_code")
	} else if !theaterReasonCodes[finding.ReasonCode] {
		v.add("theater finding has unknown reason_code %q", finding.ReasonCode)
	}
	for _, ref := range finding.TriggerEvidenceRefs {
		v.validateEvidenceRef("theater finding "+finding.ReasonCode, StatePartial, ref)
	}
}

// Once any theater finding exists, the theater row cannot remain pass because
// the row now represents an unresolved or failed anti-theater assessment.
func (v *bundleValidator) validateTheaterState() {
	row := v.rows["PC-THEATER"]
	if len(v.bundle.Packet.TheaterFindings) == 0 {
		return
	}
	if row.State == StatePass {
		v.add("PC-THEATER cannot be pass when theater findings are present")
	}
	if !theaterFindingState(row.State) {
		v.add("PC-THEATER with theater findings must be partial, fail, or cannot_verify")
	}
}

func theaterFindingState(state string) bool {
	return state == StatePartial || state == StateFail || state == StateCannotVerify
}
