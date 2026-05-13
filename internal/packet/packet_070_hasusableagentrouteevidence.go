package packet

func (c *demoFirstPacketChecker) hasUsableAgentRouteEvidence(refs []string) bool {
	// hasUsableAgentRouteEvidence keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	for _, ref := range refs {
		if c.usableAgentRouteEntry(ref) {

			return true
		}
	}
	return false
}
