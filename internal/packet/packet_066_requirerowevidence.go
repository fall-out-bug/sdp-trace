package packet

func (c *demoFirstPacketChecker) requireRowEvidence(rowID string) {
	// requireRowEvidence keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	row := c.rows[rowID]
	if len(row.EvidenceRefs) == 0 {
		c.add("demo first-packet gate requires %s retained evidence refs", rowID)
		return
	}

	c.requireUsableRowEvidenceRefs(rowID, row.EvidenceRefs)
}
