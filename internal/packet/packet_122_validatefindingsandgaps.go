package packet

func (v *bundleValidator) validateFindingsAndGaps() {
	// validateFindingsAndGaps keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	v.validateTheaterState()
	v.validateDecisionOwners()
	for _, finding := range v.bundle.Packet.TheaterFindings {
		v.validateTheaterFinding(finding)
	}
	for _, gap := range v.bundle.Packet.ResidualGaps {
		v.validateResidualGap(gap)
	}
}
