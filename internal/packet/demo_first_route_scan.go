package packet

func (c *demoFirstPacketChecker) hasUsableAgentRouteEvidence(refs []string) bool {
	for _, ref := range refs {
		if c.usableAgentRouteEntry(ref) {
			return true
		}
	}
	return false
}
