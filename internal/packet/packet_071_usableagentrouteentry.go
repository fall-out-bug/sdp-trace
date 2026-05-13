package packet

func (c *demoFirstPacketChecker) usableAgentRouteEntry(ref string) bool {
	entry := c.entryByRef[ref]
	return entry.SourceClass == "harness" && demoUsableEntry(entry, c.now) && demoRouteEvidenceObservedOpenCodeGSDMiniMax(entry)
}
