package packet

func (v *bundleValidator) validateContradictionState(rowID string, row Row) {
	if row.State != StatePartial {
		v.add("%s has contradictory evidence but state is %s, want partial", rowID, row.State)
	}
}

func (v *bundleValidator) validateContradictionGap(rowID string) {
	if !gapForRow(v.bundle.Packet.ResidualGaps, rowID) {
		v.add("%s contradictory evidence requires residual gap explanation", rowID)
	}
}
