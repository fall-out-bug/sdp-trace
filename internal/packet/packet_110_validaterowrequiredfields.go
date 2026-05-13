package packet

func (v *bundleValidator) validateRowRequiredFields(row Row) {
	v.validateRowState(row)
	v.validateRowSummary(row)
	v.validateRowOwner(row)
}
