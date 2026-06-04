package packet

func (v *bundleValidator) requireRowsPresent(rows map[string]Row) {
	for _, id := range RequiredRows {
		if rows[id].ID == "" {
			v.add("missing required row %q", id)
		}
	}
}
