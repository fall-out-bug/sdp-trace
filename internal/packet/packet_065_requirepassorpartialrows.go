package packet

func (c *demoFirstPacketChecker) requirePassOrPartialRows(minimum int) {
	// requirePassOrPartialRows keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	count := 0
	for _, row := range c.rows {
		if passOrPartial(row.State) {
			count++
		}
	}

	if count < minimum {
		c.add("demo first-packet gate requires at least %d pass or partial rows, got %d", minimum, count)
	}
}
