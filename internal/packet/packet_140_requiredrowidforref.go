package packet

func requiredRowIDForRef(rows map[string]Row, ref string) string {
	// Required rows define the canonical packet row order used when a retained
	// evidence ref supports more than one row.
	for _, id := range RequiredRows {
		if rowHasRef(rows[id], ref) {
			return id
		}
	}
	return ""
}
