package packet

func (c *demoFirstPacketChecker) requireRowEvidence(rowID string) {
	row := c.rows[rowID]
	if len(row.EvidenceRefs) == 0 {
		c.add("demo first-packet gate requires %s retained evidence refs", rowID)
		return
	}

	c.requireUsableRowEvidenceRefs(rowID, row.EvidenceRefs)
}

func (c *demoFirstPacketChecker) requireUsableRowEvidenceRefs(rowID string, refs []string) {
	for _, ref := range refs {
		entry, ok := c.entryByRef[ref]
		if !ok {
			// Base validation owns missing-ref diagnostics; this gate only adds
			// stricter demo usability checks for retained refs it can inspect.
			continue
		}
		if !demoUsableEntry(entry, c.now) {
			c.add("demo first-packet gate requires %s evidence ref %q to be retained and usable", rowID, ref)
		}
	}
}
