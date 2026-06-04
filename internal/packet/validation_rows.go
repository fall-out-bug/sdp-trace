package packet

func (v *bundleValidator) validateRows() {
	rows := map[string]Row{}
	v.indexRows(rows)
	v.requireRowsPresent(rows)
	v.rows = rows

	v.validateContradictions(rows)
	v.validateResidualCoverage(rows)
}
