package packet

func (c *demoFirstPacketChecker) requireVerificationOrReviewAssessed() {
	// requireVerificationOrReviewAssessed keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if rowAssessed(c.rows["PC-VERIFICATION"]) || rowAssessed(c.rows["PC-REVIEW"]) {
		return
	}

	c.add("demo first-packet gate requires PC-VERIFICATION or PC-REVIEW to be pass, partial, or fail")
}
