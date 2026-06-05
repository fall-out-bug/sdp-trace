package packet

func (c *demoFirstPacketChecker) requireCannotVerifyClosureCap() {
	unclosed := c.cannotVerifyRowsWithoutClosure()
	if unclosed > 1 {
		c.add("demo first-packet gate allows at most one cannot_verify row without closure path, got %d", unclosed)
	}
}

func (c *demoFirstPacketChecker) cannotVerifyRowsWithoutClosure() int {
	unclosed := 0
	for _, row := range c.rows {
		if row.State != StateCannotVerify {
			continue
		}

		// cannot_verify is allowed only when the packet names concrete closure
		// evidence; otherwise the demo gate would normalize an open trust gap.
		if !gapForRowWithClosure(c.bundle.Packet.ResidualGaps, row.ID) {
			unclosed++
		}
	}
	return unclosed
}
