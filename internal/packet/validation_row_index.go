package packet

func (v *bundleValidator) indexRows(rows map[string]Row) {
	for _, row := range v.bundle.Packet.Rows {
		if v.validateRowID(row.ID, rows) {
			rows[row.ID] = row
			v.validateRow(row)
		}
	}
}
