package packet

func (c *demoFirstPacketChecker) agentRouteEvidenceRefs() ([]string, bool) {
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
