package packet

func (c *demoFirstPacketChecker) requireAgentRouteEvidence() {
	// requireAgentRouteEvidence keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	refs, ok := c.agentRouteEvidenceRefs()
	if !ok {
		return
	}
	if c.hasUsableAgentRouteEvidence(refs) {
		return
	}
	c.add("demo first-packet gate requires PC-AGENT-ROUTE evidence from retained structured OpenCode/GSD/MiniMax harness route observation")
}
