package packet

func (v *bundleValidator) indexRows(rows map[string]Row) {
	// indexRows keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	for _, row := range v.bundle.Packet.Rows {
		if v.validateRowID(row.ID, rows) {

			rows[row.ID] = row
			v.validateRow(row)
		}
	}
}
