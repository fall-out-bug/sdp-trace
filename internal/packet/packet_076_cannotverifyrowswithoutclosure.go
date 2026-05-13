package packet

func (c *demoFirstPacketChecker) cannotVerifyRowsWithoutClosure() int {
	// cannotVerifyRowsWithoutClosure keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	unclosed := 0
	for _, row := range c.rows {
		if row.State != StateCannotVerify {
			continue
		}

		if !gapForRowWithClosure(c.bundle.Packet.ResidualGaps, row.ID) {
			unclosed++
		}
	}
	return unclosed
}
