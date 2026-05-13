package packet

func extensionRowIDForRef(rows map[string]Row, ref string) string {
	// Extension rows are scanned only after the fixed packet contract rows so
	// custom rows cannot steal contradiction attribution from required rows.
	for _, id := range sortedRowIDs(rows) {
		if !requiredRow(id) && rowHasRef(rows[id], ref) {
			return id
		}
	}
	return ""
}
