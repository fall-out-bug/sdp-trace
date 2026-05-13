package packet

func (v *bundleValidator) validateTheaterState() {
	// validateTheaterState keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
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
