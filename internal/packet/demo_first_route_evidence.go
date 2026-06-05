package packet

func (c *demoFirstPacketChecker) requireAgentRouteEvidence() {
	refs, ok := c.agentRouteEvidenceRefs()
	if !ok {
		return
	}
	if c.hasUsableAgentRouteEvidence(refs) {
		return
	}
	c.add("demo first-packet gate requires PC-AGENT-ROUTE evidence from retained structured OpenCode/GSD/MiniMax harness route observation")
}
