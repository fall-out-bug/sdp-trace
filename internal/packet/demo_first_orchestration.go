package packet

func (c *demoFirstPacketChecker) validate() Validation {
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
