package packet

func (c *demoFirstPacketChecker) usableAgentRouteEntry(ref string) bool {
	entry := c.entryByRef[ref]
	return entry.SourceClass == "harness" && demoUsableEntry(entry, c.now) && demoRouteEvidenceObservedOpenCodeGSDMiniMax(entry)
}

func demoRouteEvidenceObservedOpenCodeGSDMiniMax(entry BundleEntry) bool {
	if entry.EvidenceKind != "harness_route_observation" || syntheticEntryDigest(entry) {
		return false
	}

	return hasOpenCodeGSDMiniMax(entry.ObservedComponents)
}
