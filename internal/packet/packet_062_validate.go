package packet

func (c *demoFirstPacketChecker) validate() Validation {
	// validate keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	c.index()

	c.requireToolGenerated()
	c.requirePassOrPartialRows(4)
	c.requireRowEvidence("PC-CHANGE")
	c.requireRowEvidence("PC-MUTATION")
	c.requireAgentRouteEvidence()
	c.requireVerificationOrReviewAssessed()
	c.requireCannotVerifyClosureCap()
	state := StatePass
	if len(c.errors) > 0 {
		state = StateFail
	}
	return Validation{State: state, Errors: c.errors}
}
