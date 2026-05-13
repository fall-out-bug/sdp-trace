package packet

func (c *demoFirstPacketChecker) requireCannotVerifyClosureCap() {
	// requireCannotVerifyClosureCap keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	unclosed := c.cannotVerifyRowsWithoutClosure()
	if unclosed > 1 {
		c.add("demo first-packet gate allows at most one cannot_verify row without closure path, got %d", unclosed)
	}
}
