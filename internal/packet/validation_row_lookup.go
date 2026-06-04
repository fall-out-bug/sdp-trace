package packet

// Required rows define the canonical packet contract. Extension rows can add
// evidence but must not steal attribution from required rows.
func requiredRow(id string) bool {
	for _, required := range RequiredRows {
		if id == required {
			return true
		}
	}
	return false
}

// Ref attribution checks required rows first, then falls back to custom rows in
// deterministic order so contradictory evidence lands on the same row every run.
func rowIDForRef(rows map[string]Row, ref string) string {
	if id := requiredRowIDForRef(rows, ref); id != "" {
		return id
	}
	return extensionRowIDForRef(rows, ref)
}

// Required row order is part of the packet contract; the first required row
// citing a shared ref owns that evidence for contradiction attribution.
func requiredRowIDForRef(rows map[string]Row, ref string) string {
	for _, id := range RequiredRows {
		if rowHasRef(rows[id], ref) {
			return id
		}
	}
	return ""
}
