package packet

func (v *bundleValidator) validateRowID(rowID string, rows map[string]Row) bool {
	if !requiredRow(rowID) {
		v.add("unknown row id %q", rowID)
		return false
	}
	if rows[rowID].ID != "" {
		v.add("duplicate row id %q", rowID)
	}
	return true
}
