package packet

func (c *demoFirstPacketChecker) agentRouteEvidenceRefs() ([]string, bool) {
	// agentRouteEvidenceRefs keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	row := c.rows["PC-AGENT-ROUTE"]
	if !passOrPartial(row.State) {
		c.add("demo first-packet gate requires PC-AGENT-ROUTE must be pass or partial, got %s", row.State)
		return nil, false
	}

	if len(row.EvidenceRefs) == 0 {
		c.add("demo first-packet gate requires PC-AGENT-ROUTE retained evidence refs")
		return nil, false
	}
	return row.EvidenceRefs, true
}
